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
		ch := Collection.Find(ctx, filter, nil, nil)

		result := <-ch

		assert.Nil(t, result.Error)

		for _, document := range result.Result {
			for key, value := range document {
				fmt.Printf("Key: %v\n", key)
				fmt.Printf("Value type: %T\n", value)
				fmt.Printf("Value: %v\n", value)
			}
		}

		fmt.Printf("Result: %v", result.Result)
	}
}

func TestFind_WithoutCollection_ReturnsError(t *testing.T) {
	ctx := context.Background()
	Collection.currentCollection = nil
	ch := Collection.Find(ctx, []byte(`{"release":1987}`), nil, nil)

	result := <-ch

	assert.Nil(t, result.Result)
	assert.Equal(t, "No collection selected", result.Error.Error())
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

	chResult := Collection.Find(ctx, []byte(`{}`), []byte(`{"release": 1}`), nil)

	result := <-chResult

	assert.Nil(t, result.Error)

	fmt.Printf("Result: %v", result.Result)
}

func TestCount_WithFilter_ReturnsDocumentCount(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "amiga")

	command := []byte("{\"create\": \"countfilter\"}")
	Database.Execute(ctx, command)

	command = []byte(`{"insert": "countfilter", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987, "chipset": ["Agnus", "Denise", "Paula"]},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	_, err = Database.Execute(ctx, command)
	assert.Nil(t, err)

	Collection.SetCollection("countfilter")

	result, err := Collection.Count(ctx, []byte(`{"release": {"$gt": 1986} }`))

	assert.Nil(t, err)
	assert.Equal(t, int64(2), result)
}

func TestCount_WithEmptyFilter_ReturnsDocumentCount(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "amiga")

	command := []byte("{\"create\": \"count\"}")
	Database.Execute(ctx, command)

	command = []byte(`{"insert": "count", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987, "chipset": ["Agnus", "Denise", "Paula"]},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	_, err = Database.Execute(ctx, command)
	assert.Nil(t, err)

	Collection.SetCollection("count")

	filter := []byte{}
	assert.Equal(t, 0, len(filter))

	result, err := Collection.Count(ctx, filter)

	assert.Nil(t, err)
	assert.Equal(t, int64(3), result)
}

func TestEstimatedCount_ReturnsTotalCount(t *testing.T) {
	const connectionURI string = "mongodb://localhost"

	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connection.Connect(ctx, connection)
	err := <-ch
	assert.Nil(t, err)
	Database.UseDatabase(connectionURI, "amiga")

	command := []byte("{\"create\": \"estimated\"}")
	Database.Execute(ctx, command)

	command = []byte(`{"insert": "estimated", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987, "chipset": ["Agnus", "Denise", "Paula"]},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	_, err = Database.Execute(ctx, command)
	assert.Nil(t, err)

	Collection.SetCollection("estimated")

	result, err := Collection.EstimatedCount(ctx)

	assert.Nil(t, err)
	assert.Equal(t, int64(3), result)
}
