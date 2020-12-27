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

	"github.com/renkman/mongotui/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDatabaseExecutePrimitiveM(t *testing.T) {
	ctx := context.Background()
	connection := &models.Connection{Host: "localhost"}
	command := []byte("{\"listCommands\": 1}")

	Connect(ctx, connection)
	assert.Equal(t, "mongodb://localhost", connection.Uri)

	err := UseDatabase("mongodb://localhost", "admin")
	assert.Nil(t, err)

	result, err := Execute(ctx, command)
	assert.Nil(t, err)

	writeValueOrdered(result, 0)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDatabaseExecutePrimitiveD(t *testing.T) {
	ctx := context.Background()
	connection := &models.Connection{Host: "localhost"}
	command := []byte("{\"listCommands\": 1}")

	Connect(ctx, connection)
	assert.Equal(t, "mongodb://localhost", connection.Uri)

	err := UseDatabase("mongodb://localhost", "admin")
	assert.Nil(t, err)

	result, err := Execute(ctx, command)
	assert.Nil(t, err)

	writeValueOrdered(result, 0)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func writeValue(value interface{}, level int) {
	switch value.(type) {
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
