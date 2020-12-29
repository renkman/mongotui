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

package models

// Connection contains the values needed to connect to a MongoDB instance.
// It is used by the main application to bind the ui.FormWidget form data
// and pass it to mongo.Connect.
//
// If URI is set only this field is used, otherwise the connection URI is
// built from the other fields.
//
// The SaveConnection flag indicates whether the connection should be stored.
type Connection struct {
	Host           string
	Port           string
	User           string
	Password       string
	Replicaset     string
	TLS            bool
	URI            string
	SaveConnection bool
}
