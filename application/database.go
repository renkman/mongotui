// Copyright 2021 Jan Renken

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

package application

import (
	"context"
	"fmt"

	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/renkman/mongotui/settings"
	"github.com/renkman/mongotui/ui"
)

// Connect connects to the host with the passed *models.Connection and adds it to the
// database tree view if it was successful.

func connect(connection *models.Connection) {
	mongo.BuildConnectionURI(connection)
	if settings.CanStoreConnection && connection.SaveConnection {
		settings.StoreConnection(connection.Host, connection.URI)
	}

	ctx := context.Background()
	err := mongo.Connect(ctx, connection)
	if err != nil {
		message := fmt.Sprintf("Connection to %s failed:\n\n%s", connection.Host, err.Error())
		ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
		return
	}

	databases, err := mongo.GetDatabases(ctx, connection.URI)
	if err != nil {
		message := fmt.Sprintf("Getting databases of %s failed:\n\n%s", connection.Host, err.Error())
		ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
		return
	}
	databaseTree.AddDatabases(connection.Host, connection.URI, databases)
}

func updateDatabaseTree(connectionUri string, name string) []string {
	ctx := context.Background()
	mongo.UseDatabase(connectionUri, name)
	collections, err := mongo.GetCollections(ctx)
	if err != nil {
		message := fmt.Sprintf("Getting collections of database %s failed:\n\n%s", name, err.Error())
		ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
		return collections
	}
	return collections
}