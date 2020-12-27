package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/renkman/mongotui/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient interface {
	Connect(ctx context.Context, connection *models.Connection) error
	GetDatabases(ctx context.Context) (string, error)
}

const defaultHost string = "localhost"

var clients map[string]*mongo.Client = make(map[string]*mongo.Client)

func Connect(ctx context.Context, connection *models.Connection) error {

	buildConnectionUri(connection)

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection.Uri).SetConnectTimeout(10*time.Second))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	clients[connection.Uri] = client
	return nil
}

func GetDatabases(ctx context.Context, connectionUri string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := getClient(connectionUri)
	if err != nil {
		return []string{}, err
	}

	databases, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return databases, nil
}

func Disconnect(ctx context.Context, connectionUri string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := getClient(connectionUri)
	if err != nil {
		return err
	}
	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

func DisconnectAll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for _, client := range clients {
		client.Disconnect(ctx)
	}
	return nil
}

func buildConnectionUri(connection *models.Connection) {
	host := connection.Host
	if host == "" {
		host = defaultHost
	}
	port := ""
	if connection.Port != "" {
		port = fmt.Sprintf(":%s", connection.Port)
	}
	connection.Uri = fmt.Sprintf("mongodb://%s%s", host, port)
}

func getClient(connectionUri string) (*mongo.Client, error) {
	if client, ok := clients[connectionUri]; ok {
		return client, nil
	}
	return nil, errors.New(fmt.Sprintf("Not connected to %s", connectionUri))
}
