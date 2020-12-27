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
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var currentDatabase *mongo.Database

func UseDatabase(connectionUrl string, name string) error {
	client, err := getClient(connectionUrl)
	if err != nil {
		return err
	}
	currentDatabase = client.Database(name)
	return nil
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
	if currentDatabase == nil {
		return nil, errors.New("No database selected")
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
