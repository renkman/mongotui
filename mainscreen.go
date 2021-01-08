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

package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/renkman/mongotui/settings"
	"github.com/renkman/mongotui/ui"
	"github.com/rivo/tview"
)

func createMainSreen(ctx context.Context, app *tview.Application, pages *tview.Pages, connection *models.Connection) {
	databaseTree := ui.CreateDatabaseTreeWidget(app, pages, func(connectionUri string, name string) []string {
		mongo.UseDatabase(connectionUri, name)
		collections, err := mongo.GetCollections(ctx)
		if err != nil {
			message := fmt.Sprintf("Getting collections of database %s failed:\n\n%s", name, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return collections
		}
		return collections
	})

	resultView := ui.CreateResultTreeWidget(app, pages)
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
		databaseTree.UpdateCollections()
	})

	commandsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	commandsView.SetBorder(true).
		SetTitle("Commands")

	for i, command := range settings.GetCommands() {
		seperator := ""
		if (i+1)%5 == 0 {
			seperator = "\n"
		}
		fmt.Fprintf(commandsView, "%s%s", command.Description, seperator)
	}

	resultView.SetBorder(true).SetTitle("Result")
	editor.SetBorder(true).SetTitle("Editor")

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
		AddText("Copyright 2021 Jan Renken", true, tview.AlignRight, tcell.ColorGreenYellow)
	pages.AddPage("frame", frame, true, true)

	quitModal := ui.CreateQuitModalWidget(app, pages)

	connect := func(connection *models.Connection) {
		mongo.BuildConnectionURI(connection)
		if settings.CanStoreConnection && connection.SaveConnection {
			settings.StoreConnection(connection.Host, connection.URI)
		}

		err := mongo.Connect(ctx, connection)
		if err != nil {
			message := fmt.Sprintf("Connection to %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}

		databases, err := mongo.GetDatabases(ctx, connection.URI)
		if err != nil {
			message := fmt.Sprintf("Getting databases of %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}
		databaseTree.AddDatabases(connection.Host, connection.URI, databases)
	}

	connectionForm := ui.CreateConnectionFormWidget(app, pages, connect, settings.CanStoreConnection, settings.GetConnections, settings.GetConnectionURI)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quitModal.HandleEvent(event)
		connectionForm.HandleEvent(event)
		databaseTree.HandleEvent(event)
		resultView.HandleEvent(event)

		if event.Key() == tcell.KeyCtrlE {
			app.SetFocus(editor)
		}

		if event.Key() == tcell.KeyCtrlU {
			connectionURI := databaseTree.GetSelectedConnection()
			if connectionURI == "" {
				message := fmt.Sprintf("No host selected")
				ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
				return event
			}

			ui.CreateDatabaseFormWidget(app, pages, func(name string) {
				err := mongo.UseDatabase(connectionURI, name)
				if err != nil {
					message := fmt.Sprintf("Use database on %s failed:\n\n%s", connection.Host, err.Error())
					ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
				}
				databaseTree.AddDatabase(name)
			}).SetFocus(app)

			return event
		}

		err := databaseTree.HandleDiconnectionEvent(event, func(key string) error {
			return mongo.Disconnect(ctx, key)
		})
		if err != nil {
			message := fmt.Sprintf("Attempt to disconnect failed:\n\n%s", err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return event
		}

		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyNUL, ' ', tcell.ModNone)
		}
		return event
	})

	if connection.URI != "" {
		connect(connection)
	}
	app.SetFocus(editor)
}
