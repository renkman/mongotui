// Copyright 2021 Jan Renken

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

import "github.com/gdamore/tcell"

// Command defines the command shortcuts and the command description
type Command struct {
	Key         tcell.Key
	Description string
}

var commands []Command = []Command{
	{tcell.KeyCtrlC, "[white]Ctrl - C[darkcyan]onnect to database\t"},
	{tcell.KeyCtrlD, "[white]Ctrl - D[darkcyan]atabase tree\t"},
	{tcell.KeyCtrlE, "[white]Ctrl - E[darkcyan]nter command\t"},
	{tcell.KeyCtrlR, "[white]Ctrl - R[darkcyan]esult view\t"},
	{tcell.KeyCtrlT, "[white]Ctrl - T[darkcyan]erminate selected connection\t"},
	{tcell.KeyCtrlU, "[white]Ctrl - U[darkcyan]se database\t"},
	{tcell.KeyCtrlQ, "[white]Ctrl - Q[darkcyan]uit"}}

// GetCommands returns the defined application commands
func GetCommands() []Command {
	return commands
}
