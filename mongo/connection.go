// Copyright 2020 Jan Renken

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
	"regexp"

	"strings"
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

var (
	clients               map[string]*mongo.Client = make(map[string]*mongo.Client)
	connectionNamePattern *regexp.Regexp           = regexp.MustCompile(`mongodb(?:\+srv)*://(?:([^:]+):(?:[^@]+@)){0,1}(.*)`)
)

// Connect establishes a connection to the MongoDB instance specified by
// the passed models.Conenction and stores the resulting client in the internal
// client map with its URI as key.
func Connect(ctx context.Context, connection *models.Connection) chan error {
	BuildConnectionURI(connection)
	ch := make(chan error)

	go func() {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection.URI).SetConnectTimeout(10*time.Second))
		if err != nil {
			ch <- err
			return
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			ch <- err
			return
		}

		clients[connection.URI] = client
		ch <- nil
	}()
	return ch
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

// BuildConnectionURI builds the connection URI and adds it to the passed model
func BuildConnectionURI(connection *models.Connection) {
	if connection.URI != "" {
		connection.Host = generateClientName(connection.URI)
		return
	}
	host := connection.Host
	if host == "" {
		host = defaultHost
	}
	port := ""
	if connection.Port != "" {
		port = fmt.Sprintf(":%s", connection.Port)
	}

	credentials := ""
	if connection.User != "" {
		credentials = fmt.Sprintf("%s:%s@", connection.User, connection.Password)
	}

	var options []string
	if connection.Replicaset != "" {
		options = append(options, fmt.Sprintf("replicaSet=%s", connection.Replicaset))
	}
	if connection.TLS {
		options = append(options, "tls=true")
	}

	optionParameters := strings.Join(options, "&")
	if optionParameters != "" {
		optionParameters = fmt.Sprintf("?%s", optionParameters)
	}

	connection.URI = fmt.Sprintf("mongodb://%s%s%s%s", credentials, host, port, optionParameters)
}

func getClient(connectionURI string) (*mongo.Client, error) {
	if client, ok := clients[connectionURI]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("Not connected to %s", connectionURI)
}

func generateClientName(connectionURI string) string {
	result := connectionNamePattern.FindStringSubmatch(connectionURI)

	if len(result) == 0 || result[2] == "" {
		return connectionURI
	}
	if result[1] == "" {
		return result[2]
	}
	return strings.Join(result[1:], "@")
}
