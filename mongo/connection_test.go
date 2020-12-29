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

func Test_BuildConnectionURI_WithHost(t *testing.T) {
	model := &models.Connection{Host: "localhost"}
	BuildConnectionURI(model)

	assert.Equal(t, "mongodb://localhost", model.URI)
}

func Test_BuildConnectionURI_WithURIWithoutCredentials(t *testing.T) {
	model := &models.Connection{URI: "mongodb://localhost"}
	BuildConnectionURI(model)

	assert.Equal(t, "localhost", model.Host)
}

func Test_BuildConnectionURI_WithURIWithCredentials(t *testing.T) {
	model := &models.Connection{URI: "mongodb://foo:bar@localhost"}
	BuildConnectionURI(model)

	assert.Equal(t, "foo@localhost", model.Host)
}

func Test_BuildConnectionURI_WithURIWithProtocolOnly(t *testing.T) {
	model := &models.Connection{URI: "mongodb://"}
	BuildConnectionURI(model)

	assert.Equal(t, "mongodb://", model.Host)
}

func Test_BuildConnectionURI_WithInvalidURI(t *testing.T) {
	model := &models.Connection{URI: "foobar"}
	BuildConnectionURI(model)

	assert.Equal(t, "foobar", model.Host)
}
