package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collection struct {
	currentCollection *mongo.Collection
}

var Collection *collection = &collection{}

func (collection *collection) SetCollection(name string) {
	collection.currentCollection = Database.currentDatabase.Collection(name)
}

func (collection *collection) Find(ctx context.Context, filter []byte, sort []byte, project []byte) ([]map[string]interface{}, error) {
	if collection.currentCollection == nil {
		return nil, fmt.Errorf("No collection selected")
	}

	if filter == nil || len(filter) == 0 {
		filter = []byte(`{}`)
	}

	if sort == nil || len(sort) == 0 {
		sort = []byte(`{}`)
	}

	if project == nil || len(project) == 0 {
		project = []byte(`{}`)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filterBson, err := unmarshal(filter)
	if err != nil {
		return nil, err
	}

	sortBson, err := unmarshal(sort)
	if err != nil {
		return nil, err
	}

	projectBson, err := unmarshal(project)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(sortBson).SetProjection(projectBson)

	cursor, err := collection.currentCollection.Find(ctx, filterBson, opts)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func unmarshal(command []byte) (interface{}, error) {
	if command == nil {
		return nil, nil
	}
	var commandBson interface{}
	err := bson.UnmarshalExtJSON(command, true, &commandBson)
	if err != nil {
		return nil, err
	}
	return commandBson, nil
}
