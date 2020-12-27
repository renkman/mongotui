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

package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/renkman/mongotui/ui"
	"github.com/rivo/tview"
)

func CreateMainSreen(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	resultView := ui.CreateResultTree(app, pages)
	editor := tview.NewInputField().
		SetLabel("Command: ").
		SetFieldWidth(200)
	editor.SetDoneFunc(func(key tcell.Key) {
		result, err := mongo.Execute(ctx, []byte(editor.GetText()))
		if err != nil {
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, err.Error())
			return
		}
		resultView.SetResult(result)
	})

	commandsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter)

	commandsView.SetBorder(true).
		SetTitle("Commands")

	fmt.Fprint(commandsView, "\n[white]Ctrl - Q [darkcyan]Quit\t[white]Ctrl - C[darkcyan] Connect to database")

	resultView.SetBorder(true).SetTitle("Result")
	editor.SetBorder(true).SetTitle("Editor")

	databaseTree := ui.GetDatabaseTree(app, pages, func(connectionUri string, name string) []string {
		mongo.UseDatabase(connectionUri, name)
		collections, err := mongo.GetCollections(ctx)
		if err != nil {
			message := fmt.Sprintf("Getting collections of database %s failed:\n\n%s", name, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return collections
		}
		return collections
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandsView, 5, 1, false).
		AddItem(tview.NewFlex().
			AddItem(databaseTree, 40, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(editor, 3, 1, false).
				AddItem(resultView, 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Statistics"), 10, 1, false),
				0, 1, false),
			0, 1, false)

	frame := tview.NewFrame(flex).
		AddText("MongoTUI - MongoDB crawler", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Copyright 2020 Jan Renken", true, tview.AlignRight, tcell.ColorYellow)
	pages.AddPage("frame", frame, true, true)

	quitModal := ui.GetQuitModalWidget(app, pages)

	connectionForm := ui.GetConnectionFormWidget(app, pages, func(connection *models.Connection) {
		err := mongo.Connect(ctx, connection)
		if err != nil {
			message := fmt.Sprintf("Connection to %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}

		databases, err := mongo.GetDatabases(ctx, connection.Uri)
		if err != nil {
			message := fmt.Sprintf("Getting databeses of %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}
		databaseTree.AddDatabases(connection.Host, connection.Uri, databases)
		pages.RemovePage("connection")
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quitModal.SetEvent(event)
		connectionForm.SetEvent(event)
		databaseTree.SetEvent(event)

		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyNUL, ' ', tcell.ModNone)
		}
		return event
	})

	app.SetFocus(editor)
}
