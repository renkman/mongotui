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
	"flag"

	"github.com/99designs/keyring"
	"github.com/renkman/mongotui/models"
)

// InitCommandLineArgs initializes the command line arguments. Currently, it is just
// -c to set a connection URI to connect to a MongoDB instance directly after
// application start.
func InitCommandLineArgs(connection *models.Connection) {
	flag.StringVar(&connection.URI,
		"c",
		"",
		"MongoDB Connection URI to connect directly after the application start")
}

func checkKeyring() bool {
	backends := keyring.AvailableBackends()
	return len(backends) > 0
}
