package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseError struct {
	message string
}

var currentDatabase *mongo.Database

func UseDatabase(name string) {
	currentDatabase = currentClient.Database(name)
}

func GetCollections(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	collections, err := currentDatabase.ListCollectionNames(ctx, currentDatabase)
	if err != nil {
		return []string{}, err
	}
	return collections, nil
}

func Execute(ctx context.Context, command []byte) (interface{}, error) {
	if currentClient == nil {
		return nil, &DatabaseError{"Not connected"}
	}

	if currentDatabase == nil {
		return nil, &DatabaseError{"No database selected"}
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var commandBson interface{}
	err := bson.UnmarshalExtJSON(command, true, &commandBson)
	if err != nil {
		return nil, err
	}

	opts := options.RunCmd().SetReadPreference(readpref.Primary())

	var result interface{}
	err = currentDatabase.RunCommand(ctx, commandBson, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf(e.message)
}
