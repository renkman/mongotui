// This file is part of MongoTUI.

// MongoTUI is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// MongoTUI is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with MongoTUI.  If not, see <http://www.gnu.org/licenses/>.

package mongo

import (
	"context"
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

// Connect establishes a connection to the MongoDB instance specified by
// the passed models.Conenction and stores the resulting client in the internal
// client map with its URI as key.
func Connect(ctx context.Context, connection *models.Connection) error {
	if connection.URI == "" {
		buildConnectionURI(connection)
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection.URI).SetConnectTimeout(10*time.Second))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	clients[connection.URI] = client
	return nil
}

// GetDatabases returns the databases of the MongoDB instance specified by the
// passed connectionURI as string slice.
func GetDatabases(ctx context.Context, connectionURI string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := getClient(connectionURI)
	if err != nil {
		return []string{}, err
	}

	databases, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return []string{}, err
	}
	return databases, nil
}

// Disconnect disconnects from the MongoDB instance specified by
// the passed connectionURI and removes the related entry from the
// internal clients map.
func Disconnect(ctx context.Context, connectionURI string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := getClient(connectionURI)
	if err != nil {
		return err
	}
	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	delete(clients, connectionURI)

	return nil
}

// DisconnectAll disconnects from all connected MongoDB instances and
// cleans up the internal clients map.
func DisconnectAll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for key, client := range clients {
		client.Disconnect(ctx)
		delete(clients, key)
	}
	return nil
}

func buildConnectionURI(connection *models.Connection) {
	host := connection.Host
	if host == "" {
		host = defaultHost
	}
	port := ""
	if connection.Port != "" {
		port = fmt.Sprintf(":%s", connection.Port)
	}
	connection.URI = fmt.Sprintf("mongodb://%s%s", host, port)
}

func getClient(connectionURI string) (*mongo.Client, error) {
	if client, ok := clients[connectionURI]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("Not connected to %s", connectionURI)
}
