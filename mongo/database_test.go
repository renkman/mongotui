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

func TestDatabaseExecute(t *testing.T) {
	ctx := context.Background()
	connection := models.Connection{Host: "localhost"}
	command := []byte("{\"listCommands\": 1}")

	Connect(ctx, connection)
	UseDatabase("admin")
	result, err := Execute(ctx, command)

	writeValue(result, 0)

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
