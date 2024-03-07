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
			{"_id":3, "name":"Amiga 4000", "release": 1992},
			{"_id":4, "object": {"foo":"Bar"}}
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

// Key: _id
// Value type: int32
// Value: 1
// Key: name
// Value type: string
// Value: Amiga 500
// Key: release
// Value type: int32
// Value: 1987
// Key: chipset
// Value type: primitive.A
// Value: [Agnus Denise Paula]
// Key: _id
// Value type: int32
// Value: 2
// Key: name
// Value type: string
// Value: Amiga 1000
// Key: release
// Value type: int32
// Value: 1985
// Key: release
// Value type: int32
// Value: 1992
// Key: _id
// Value type: int32
// Value: 3
// Key: name
// Value type: string
// Value: Amiga 4000
// Key: _id
// Value type: int32
// Value: 4
// Key: object
// Value type: map[string]interface {}
// Value: map[foo:Bar]
// Result: [map[_id:1 chipset:[Agnus Denise Paula] name:Amiga 500 release:1987] map[_id:2 name:Amiga 1000 release:1985] map[_id:3 name:Amiga 4000 release:1992] map[_id:4 object:map[foo:Bar]]]--- PASS: TestFind_ReturnsResult (0.02s)
