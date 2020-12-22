package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/renkman/mongotui/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoClient interface {
	Connect(ctx context.Context, connection models.Connection) error
	GetDatabases(ctx context.Context) (string, error)
}

const defaultPort string = "27017"

var currentClient *mongo.Client

func Connect(ctx context.Context, connection models.Connection) error {

	port := connection.Port
	if port == "" {
		port = defaultPort
	}
	uri := fmt.Sprintf("mongodb://%s:%s", connection.Host, port)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	currentClient = client
	return nil
}

func GetDatabases(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	databases, err := currentClient.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return databases, nil
}

func Disconnect(ctx context.Context) error {
	if currentClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := currentClient.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
