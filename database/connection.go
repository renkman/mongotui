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

package database

import (
	"context"

	"github.com/renkman/mongotui/models"
)

// Connecter provides an interface for database connection functions.
type Connecter interface {
	Connect(ctx context.Context, connection *models.Connection) chan error
	Disconnect(ctx context.Context, connectionURI string) error
	DisconnectAll(ctx context.Context) error
	GetDatabases(ctx context.Context, connectionURI string) ([]string, error)
}
