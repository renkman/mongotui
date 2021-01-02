package settings

import (
	"fmt"
	//"os"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
)

func TestKeyring_Get_Set(t *testing.T) {
	// if os.Getenv("AGENT_ID") != "" {
	// 	t.Skip("Keyring tests do not run on CI environment")
	// }

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})
	assert.Nil(t, err)

	err = ring.Set(keyring.Item{
		Key:   "connection",
		Data:  []byte("secret mongo connection"),
		Label: "MongoDB Connection",
	})
	assert.Nil(t, err)

	item, err := ring.Get("connection")
	assert.Nil(t, err)

	assert.Equal(t, []byte("secret mongo connection"), item.Data)
}

func TestKeyring_Keys(t *testing.T) {
	// if os.Getenv("AGENT_ID") != "" {
	// 	t.Skip("Keyring tests do not run on CI environment")
	// }

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})
	assert.Nil(t, err)

	keys, err := ring.Keys()
	assert.Nil(t, err)

	fmt.Printf("%v", keys)
}
