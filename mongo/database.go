package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var currentDatabase *mongo.Database

func UseDatabase(name string) {
	currentDatabase = currentClient.Database(name)
}

func GetCollections(foo context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(foo, 10*time.Second)
	defer cancel()

	collections, err := currentDatabase.ListCollectionNames(ctx, currentDatabase)
	if err != nil {
		return []string{}, err
	}
	return collections, nil
}
