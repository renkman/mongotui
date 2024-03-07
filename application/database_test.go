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
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

type testConnection struct {
	CurrentDatabase string
}

type testDatabase struct{}

const (
	timeout            time.Duration = 1 * time.Second
	validHost          string        = "mongodb://commodore.com"
	unreachableHost    string        = "mongodb://apple.com"
	successfulDatabase string        = "homecomputers"
	failedDatabase     string        = "iCrap"
)

var connecter *testConnection = &testConnection{"homecomputers"}

// Test connection mocks
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

// Test database mocks
func (t *testDatabase) UseDatabase(connectionURI string, name string) error {
	if connectionURI == validHost && name == successfulDatabase {
		return nil
	}
	return fmt.Errorf("Using database failed")
}

func (t *testDatabase) GetCollections(ctx context.Context) ([]string, error) {
	return []string{"Amiga 500", "Amiga 1000"}, nil
}

func (t *testDatabase) Drop(ctx context.Context) error {
	return fmt.Errorf("Not implemented yet")
}

func (t *testDatabase) Execute(ctx context.Context, command []byte) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (t *testDatabase) GetCurrentDatabaseName() (string, error) {
	return "", fmt.Errorf("Not implemented yet")
}

// Tests
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

	databases := getDatabases(t)
	assert.Len(t, databases, 1)
	assert.Equal(t, connecter.CurrentDatabase, databases[0].GetText())
}

func TestUpdateDatabaseTree_updatesDatabase(t *testing.T) {
	Database = &testDatabase{}

	collections := updateDatabaseTree(validHost, "Amiga")
	assert.Len(t, collections, 2)
}

func getDatabases(t *testing.T) []*tview.TreeNode {
	t.Cleanup(func() {
		databaseTree.GetRoot().ClearChildren()
	})

	clients := databaseTree.GetRoot().GetChildren()
	assert.Len(t, clients, 1)
	assert.Equal(t, validHost, clients[0].GetText())

	return clients[0].GetChildren()
}
