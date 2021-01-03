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

/*
MongoTUI is a MongoDb TUI client which allows to connect to multiple MongoDB
instances.

	Usage of mongotui:
	  -c string
	        MongoDB Connection URI to connect directly after the application start

Press <Ctrl>-<c> to connect to a MongoDB instance, you can enter the connection parameters individually or the connection URI as well. Notice that the connection URI always wins, if the individual fields and the connection URI are filled.

The open database connections, accessed by <Ctrl>-<d>, their databases and collections are displayed as a tree view in the left application panel. You can navigate through the nodes with the arrow keys or left-click on them.

The command editor is accessible by <Ctrl>-<e>. The commands are fired on the database which are selected in the tree view by pressing <Enter> or <Return> in the command editor.

The command result is shown in the result panel as a tree view. You can access it with <Ctrl>-<r> and navigate through the nodes with the arrow keys.

<Ctrl>-<t> disconnects (terminates) the open connections.

<Ctrl>-<q> disconnects the open connections and quits the application.
*/
package main
