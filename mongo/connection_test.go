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
	"testing"

	"github.com/renkman/mongotui/models"
	"github.com/stretchr/testify/assert"
)

var connectionTestData = []struct {
	in  *models.Connection
	out string
}{
	{&models.Connection{}, "mongodb://localhost"},
	{&models.Connection{Host: "commodore.com"}, "mongodb://commodore.com"},
	{&models.Connection{Host: "commodore.com", User: "Jay", Password: "Miner"}, "mongodb://Jay:Miner@commodore.com"},
	{&models.Connection{Host: "commodore.com", Port: "27017"}, "mongodb://commodore.com:27017"},
	{&models.Connection{Host: "commodore.com", User: "Jay", Password: "Miner", Port: "27017"}, "mongodb://Jay:Miner@commodore.com:27017"},
	{&models.Connection{Host: "commodore.com", Port: "27017", Replicaset: "foo"}, "mongodb://commodore.com:27017?replicaSet=foo"},
	{&models.Connection{Host: "commodore.com", Port: "27017", Replicaset: "foo", TLS: true}, "mongodb://commodore.com:27017?replicaSet=foo&tls=true"},
	{&models.Connection{Host: "commodore.com", User: "Jay", Password: "Miner", Port: "27017", Replicaset: "foo", TLS: true}, "mongodb://Jay:Miner@commodore.com:27017?replicaSet=foo&tls=true"},
}

func Test_BuildConnectionURI_WitIndividualFields_BuildURI(t *testing.T) {
	for _, connectionTest := range connectionTestData {
		t.Run(connectionTest.out, func(t *testing.T) {
			BuildConnectionURI(connectionTest.in)

			assert.Equal(t, connectionTest.out, connectionTest.in.URI)
		})
	}
}

var uriTestData = []struct {
	in  string
	out string
}{
	{"mongodb://localhost", "localhost"},
	{"mongodb+srv://commodore.com", "commodore.com"},
	{"mongodb://foo:bar@localhost", "foo@localhost"},
	{"mongodb+srv://foo:bar@localhost", "foo@localhost"},
	{"mongodb://", "mongodb://"},
	{"foobar", "foobar"},
}

func Test_BuildConnectionURI_WithURI_SetsHost(t *testing.T) {
	for _, uriTest := range uriTestData {
		t.Run(uriTest.in, func(t *testing.T) {
			model := &models.Connection{URI: uriTest.in}
			BuildConnectionURI(model)

			assert.Equal(t, uriTest.out, model.Host)
		})
	}
}

var connectionHostTestData = []struct {
	in  *models.Connection
	out string
}{
	{&models.Connection{Host: "apple.com", URI: "mongodb://commodore.com"}, "commodore.com"},
	{&models.Connection{Host: "apple.com", URI: "mongodb+srv://commodore.com"}, "commodore.com"},
	{&models.Connection{Host: "apple.com", URI: "mongodb://jay:miner@commodore.com"}, "jay@commodore.com"},
	{&models.Connection{Host: "apple.com", URI: "mongodb+srv://jay:miner@commodore.com"}, "jay@commodore.com"},
}

func Test_BuildConnectionURI_WithURIAndHost_SetsHostToURIhostname(t *testing.T) {
	for _, connectionHostTest := range connectionHostTestData {
		t.Run(connectionHostTest.out, func(t *testing.T) {
			BuildConnectionURI(connectionHostTest.in)
			assert.Equal(t, connectionHostTest.out, connectionHostTest.in.Host)
		})
	}
}
