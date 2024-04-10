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

	filters := [][]byte{[]byte(`{}`), []byte{}, nil}

	for _, filter := range filters {
		result, err := Collection.Find(ctx, filter, nil, nil)
		assert.Nil(t, err)

		for _, document := range result {
			for key, value := range document {
				fmt.Printf("Key: %v\n", key)
				fmt.Printf("Value type: %T\n", value)
				fmt.Printf("Value: %v\n", value)
			}
		}

		fmt.Printf("Result: %v", result)
	}
}

func TestFind_WithoutCollection_ReturnsError(t *testing.T) {
	ctx := context.Background()
	Collection.currentCollection = nil
	result, err := Collection.Find(ctx, []byte(`{"release":1987}`), nil, nil)

	assert.Nil(t, result)
	assert.Equal(t, "No collection selected", err.Error())
}

func TestFind_WithOrder_ReturnsOrderedResult(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "amiga")

	command := []byte("{\"create\": \"ordered\"}")
	Database.Execute(ctx, command)

	command = []byte(`{"insert": "ordered", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987, "chipset": ["Agnus", "Denise", "Paula"]},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	_, err = Database.Execute(ctx, command)
	assert.Nil(t, err)

	Collection.SetCollection("ordered")

	result, err := Collection.Find(ctx, []byte(`{}`), []byte(`{"release": 1}`), nil)
	assert.Nil(t, err)

	fmt.Printf("Result: %v", result)
}
