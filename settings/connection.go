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
	"github.com/99designs/keyring"
)

func StoreConnection(key string, connectionURI string) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})
	if err != nil {
		return err
	}

	err = ring.Set(keyring.Item{
		Key:   key,
		Data:  []byte(connectionURI),
		Label: "MongoDB Connection",
	})
	return err
}

func GetConnections() ([]string, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})
	if err != nil {
		return nil, err
	}

	keys, err := ring.Keys()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func GetConnectionURI(key string) (string, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "mongoTUI",
	})
	if err != nil {
		return "", err
	}

	value, err := ring.Get(key)
	if err != nil {
		return "", err
	}

	return string(value.Data), nil
}
