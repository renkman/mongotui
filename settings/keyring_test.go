package settings

import (
	"fmt"
	"os"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
)

func TestKeyring_Get_Set(t *testing.T) {
	if os.Getenv("AGENT_ID") != "" {
		t.Skip("Keyring tests do not run on CI environment")
	}

	ring, _ := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})

	_ = ring.Set(keyring.Item{
		Key:   "Connection",
		Data:  []byte("secret mongo connection"),
		Label: "MongoDB Connection",
	})

	item, _ := ring.Get("foo")

	assert.Equal(t, []byte("secret-bar"), item.Data)
}

func TestKeyring_Keys(t *testing.T) {
	if os.Getenv("AGENT_ID") != "" {
		t.Skip("Keyring tests do not run on CI environment")
	}

	ring, _ := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})

	keys, _ := ring.Keys()

	fmt.Printf("%v", keys)
}
