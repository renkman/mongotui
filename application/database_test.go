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
	"testing"
	"time"

	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/stretchr/testify/assert"
)

func TestConnect_TriesConnectAndCancelsAfter30Seconds(t *testing.T) {
	draw = func() {}

	start := time.Now()
	connection := &models.Connection{Host: "foo"}
	Connect(connection)
	stop := time.Now()

	result := stop.Sub(start)

	assert.LessOrEqual(t, 30.0, result.Seconds())
}

func TestConnect_ConnectsToDatabaseAndSetsCurrentDatabase(t *testing.T) {
	const uri string = "mongodb://localhost"
	connection := &models.Connection{URI: uri}
	Connect(connection)

	result, err := mongo.GetDatabases(context.TODO(), uri)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
}
