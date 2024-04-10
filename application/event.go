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

package application

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/ui"
)

func handleEditorEvent(key tcell.Key) {
	ctx := context.Background()
	// result, err := Database.Execute(ctx, []byte(editor.GetText()))
	result, err := Collection.Find(ctx, []byte(filterEditor.GetText()), []byte(sortEditor.GetText()), []byte(projectEditor.GetText()))
	//result, err := Collection.Find(ctx, []byte(filterEditor.GetText()), nil, nil)
	if err != nil {
		ui.CreateMessageModalWidget(app, pages, ui.TypeError, err.Error())
		return
	}
	resultView.SetResult(result)
	databaseTree.UpdateCollections()
}

func handleUseDatabaseEvent(event *tcell.EventKey) *tcell.EventKey {
	connectionURI := databaseTree.GetSelectedConnection()
	if connectionURI == "" {
		message := fmt.Sprintf("No host selected")
		ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
		return event
	}

	ui.CreateDatabaseFormWidget(app, pages, func(name string) {
		err := Database.UseDatabase(connectionURI, name)
		if err != nil {
			message := fmt.Sprintf("Use database on %s failed:\n\n%s", connectionURI, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
		}
		databaseTree.AddDatabase(name)
	}).SetFocus(app)

	return event
}
