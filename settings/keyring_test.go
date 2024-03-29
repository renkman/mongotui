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

package settings

import (
	"os"
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
)

func TestKeyring_Get_Set(t *testing.T) {
	if os.Getenv("AGENT_ID") != "" || os.Getenv("CI") == "true" {
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

	keys, _ := ring.Keys()

	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "Connection", keys[0])

	item, _ := ring.Get("Connection")

	assert.Equal(t, []byte("secret mongo connection"), item.Data)
}

func TestKeyring_Keys(t *testing.T) {
	if os.Getenv("AGENT_ID") != "" || os.Getenv("CI") == "true" {
		t.Skip("Keyring tests do not run on CI environment")
	}

	ring, _ := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})

	keys, _ := ring.Keys()

	assert.Equal(t, 3, len(keys))
}
