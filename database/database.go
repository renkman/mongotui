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

import "context"

type Database interface {
	UseDatabase(connectionURI string, name string) error
	GetCollections(ctx context.Context) ([]string, error)
	Execute(ctx context.Context, command []byte) (interface{}, error)
	Drop(ctx context.Context) error
	GetCurrentDatabaseName() (string, error)
}
