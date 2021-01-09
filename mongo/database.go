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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var currentDatabase *mongo.Database

// UseDatabase sets the current database specified by the passed name of the MongoDB instance
// specified by the passed connectionURI.
// Since the MongoDB use command is used, the database will be created if it does not
// exist.
func UseDatabase(connectionURI string, name string) error {
	client, err := getClient(connectionURI)
	if err != nil {
		return err
	}
	currentDatabase = client.Database(name)
	return nil
}

// GetCollections returns the collections of the current database, which is
// set by UseDatabase.
func GetCollections(ctx context.Context) ([]string, error) {
	if currentDatabase == nil {
		return nil, fmt.Errorf("No database selected")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	collections, err := currentDatabase.ListCollectionNames(ctx, currentDatabase)
	if err != nil {
		return []string{}, err
	}
	return collections, nil
}

// Execute executes the passed command on the current database, which is set by
// UseDatabase.
func Execute(ctx context.Context, command []byte) (interface{}, error) {
	if currentDatabase == nil {
		return nil, fmt.Errorf("No database selected")
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

// Drop drops the current database.
func Drop(ctx context.Context) error {
	if currentDatabase == nil {
		return fmt.Errorf("No database selected")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := currentDatabase.Drop(ctx)
	return err
}

// GetCurrentDatabaseName returns the current database name.
func GetCurrentDatabaseName() (string, error) {
	if currentDatabase == nil {
		return "", fmt.Errorf("No database selected")
	}

	return currentDatabase.Name(), nil
}
