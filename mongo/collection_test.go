package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/renkman/mongotui/models"
	"github.com/stretchr/testify/assert"
)

func TestFind_ReturnsResult(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "homecomputer")

	command := []byte("{\"create\": \"systems\"}")
	Database.Execute(ctx, command)

	command = []byte(`{"insert": "systems", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987, "chipset": ["Agnus", "Denise", "Paula"]},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	_, err = Database.Execute(ctx, command)
	assert.Nil(t, err)

	Collection.SetCollection("systems")

	result, err := Collection.Find(ctx, []byte(`{}`), nil, nil)
	assert.Nil(t, err)

	fmt.Printf("Result: %v", result)
}

func TestFind_WithoutCollection_ReturnsError(t *testing.T) {
	ctx := context.Background()
	Collection.currentCollection = nil
	result, err := Collection.Find(ctx, []byte(`{"release":1987}`), nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, "No collection selected", err.Error())
}

func TestFind_WithoutFilter_ReturnsError(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "no_filter_test")

	command := []byte("{\"create\": \"fail\"}")
	Database.Execute(ctx, command)

	Collection.SetCollection("fail")
	result, err := Collection.Find(ctx, nil, nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, "Argument filter must not be nil", err.Error())
}
