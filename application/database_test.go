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
	"testing"
	"time"

	"github.com/renkman/mongotui/models"
	"github.com/stretchr/testify/assert"
)

type testConnection struct {
	CurrentDatabase string
}

const (
	timeout         time.Duration = 1 * time.Second
	validHost       string        = "mongodb://commodore.com"
	unreachableHost string        = "mongodb://apple.com"
)

var connecter *testConnection = &testConnection{"homecomputers"}

func (t *testConnection) Connect(ctx context.Context, connection *models.Connection) chan error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	ch := make(chan error)
	go func() {
		defer cancel()
		if connection.Host == validHost {
			ch <- nil
			return
		}
		if connection.Host == unreachableHost {
			<-ctx.Done()
			ch <- fmt.Errorf("Connection failed!")
			return
		}
	}()
	return ch
}

func (t *testConnection) Disconnect(ctx context.Context, connectionURI string) error {
	return fmt.Errorf("Not implemented yet")
}
func (t *testConnection) DisconnectAll(ctx context.Context) error {
	return fmt.Errorf("Not implemented yet")
}
func (t *testConnection) GetDatabases(ctx context.Context, connectionURI string) ([]string, error) {
	return []string{t.CurrentDatabase}, nil
}

func TestConnect_TriesConnectAndCancelsAfterTimeout(t *testing.T) {
	draw = func() {}

	start := time.Now()
	connection := &models.Connection{Host: unreachableHost}
	Connect(connecter, connection)
	stop := time.Now()

	result := stop.Sub(start)

	assert.LessOrEqual(t, timeout.Seconds(), result.Seconds())

	clients := databaseTree.TreeView.GetRoot().GetChildren()
	assert.Empty(t, clients)
}

func TestConnect_ConnectsToDatabaseAndSetsCurrentDatabase(t *testing.T) {
	connection := &models.Connection{Host: validHost}
	Connect(connecter, connection)

	clients := databaseTree.GetRoot().GetChildren()
	assert.Len(t, clients, 1)
	assert.Equal(t, validHost, clients[0].GetText())

	databases := clients[0].GetChildren()
	assert.Len(t, databases, 1)
	assert.Equal(t, connecter.CurrentDatabase, databases[0].GetText())

	databaseTree.GetRoot().ClearChildren()
}
