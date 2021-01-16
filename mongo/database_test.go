// Copyright 2020 Jan Renken

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

package mongo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	//"github.com/google/uuid"
	"github.com/renkman/mongotui/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const connectionURI string = "mongodb://localhost"

func TestExecute_WithValidCommand_ReturnsPrimitiveD(t *testing.T) {
	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	command := []byte(`{"listCommands": 1}`)

	ch := Connect(ctx, connection)
	err := <-ch
	assert.Equal(t, connectionURI, connection.URI)

	err = UseDatabase(connectionURI, "admin")
	assert.Nil(t, err)

	result, err := Execute(ctx, command)
	assert.Nil(t, err)

	assert.IsType(t, primitive.D{}, result)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestUse_WithNewDatabase_CreatesDatabase(t *testing.T) {
	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}

	ch := Connect(ctx, connection)
	err := <-ch
	assert.Equal(t, connectionURI, connection.URI)

	err = UseDatabase(connectionURI, "foobar")
	assert.Nil(t, err)

	command := []byte(`{"create": "foo"}`)
	result, err := Execute(ctx, command)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	databases, err := GetDatabases(ctx, connectionURI)
	assert.Nil(t, err)
	assert.Contains(t, databases, "foobar")
}

func TestExecute_WithInsertAndFind_ReturnsCursor(t *testing.T) {
	ctx := context.Background()
	connection := &models.Connection{URI: connectionURI}
	ch := Connect(ctx, connection)
	err := <-ch
	UseDatabase(connectionURI, "commodore")

	command := []byte("{\"create\": \"systems\"}")
	Execute(ctx, command)

	command = []byte(`{"insert": "systems", "documents": [
			{"_id":1, "name":"Amiga 500", "release": 1987},
			{"_id":2, "name":"Amiga 1000", "release": 1985},
			{"_id":3, "name":"Amiga 4000", "release": 1992}
		]}`)
	Execute(ctx, command)

	command = []byte(`{"find":"systems"}`)
	result, err := Execute(ctx, command)
	assert.Nil(t, err)

	writeValue(result, 0)
}

func writeValue(value interface{}, level int) {
	switch value.(type) {
	case primitive.A:
		resultArray := value.(primitive.A)
		for k, v := range resultArray {
			fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(v))
			for i := 0; i < level; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("%v\n", k)
			writeValueOrdered(v, level+1)
		}
	case primitive.D:
		resultMap := value.(primitive.D).Map()
		for k, v := range resultMap {
			fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(v))
			for i := 0; i < level; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("%v\n", k)
			writeValue(v, level+1)
		}
	default:
		fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(value))
		for i := 0; i < level; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("%v\n", value)
	}
}

func writeValueOrdered(value interface{}, level int) {
	switch value.(type) {
	case primitive.A:
		resultArray := value.(primitive.A)
		for k, v := range resultArray {
			fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(v))
			for i := 0; i < level; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("%v\n", k)
			writeValueOrdered(v, level+1)
		}
	case primitive.D:
		resultMap := value.(primitive.D)
		for k, v := range resultMap {
			fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(v))
			for i := 0; i < level; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("%v\n", k)
			writeValueOrdered(v, level+1)
		}
	case primitive.E:
		resultElement := value.(primitive.E)
		fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(resultElement.Value))
		for i := 0; i < level; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("%v\n", resultElement.Key)
		writeValueOrdered(resultElement.Value, level+1)
	default:
		fmt.Printf("Level: %v\t%v\t", level, reflect.TypeOf(value))
		for i := 0; i < level; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("%v\n", value)
	}
}
