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
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/database"
	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/renkman/mongotui/settings"
	"github.com/renkman/mongotui/ui"
	"github.com/rivo/tview"
)

var (
	app           *tview.Application
	pages         *tview.Pages
	databaseTree  *ui.DatabaseTreeWidget
	resultView    *ui.ResultTreeWidget
	filterEditor  *tview.InputField
	sortEditor    *tview.InputField
	projectEditor *tview.InputField
	editorView    *tview.Flex
	commandsView  *tview.TextView
	draw          func()
	getConnection func() database.Connecter
)

func init() {
	app = tview.NewApplication()
	pages = tview.NewPages()
	draw = func() { app.Draw() }
	getConnection = func() database.Connecter { return mongo.Connection }

	databaseTree = ui.CreateDatabaseTreeWidget(app, pages, updateDatabaseTree, setCollection)

	resultView = ui.CreateResultTreeWidget(app, pages)
	resultView.SetBorder(true).SetTitle("Result")

	filterEditor = tview.NewInputField().
		SetLabel("Filter ").
		SetLabelWidth(10).
		SetFieldWidth(200).
		SetDoneFunc(handleEditorEvent)

	sortEditor = tview.NewInputField().
		SetLabel("Sort ").
		SetLabelWidth(10).
		SetFieldWidth(200).
		SetDoneFunc(handleEditorEvent)

	projectEditor = tview.NewInputField().
		SetLabel("Project ").
		SetLabelWidth(10).
		SetFieldWidth(200).
		SetDoneFunc(handleEditorEvent)

	editorView = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(filterEditor, 1, 1, false).
		AddItem(sortEditor, 1, 1, false).
		AddItem(projectEditor, 1, 1, false)
	editorView.SetBorder(true).SetTitle("Editor")

	commandsView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	commandsView.SetBorder(true).
		SetTitle("Commands")
	setupCommandsView()

	quitModal := ui.CreateQuitModalWidget(app, pages)
	connectionForm := ui.CreateConnectionFormWidget(app, pages,
		func(connection *models.Connection) {
			Connect(getConnection(), connection)
		},
		settings.CanStoreConnection, settings.GetConnections, settings.GetConnectionURI)
	dropDatabaseForm := ui.CreateDropDatabaseModalWidget(app, pages, getCurrentDatabase, dropDatabase)

	buildMainScreen()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quitModal.HandleEvent(event)
		connectionForm.HandleEvent(event)
		databaseTree.HandleEvent(event)
		resultView.HandleEvent(event)
		dropDatabaseForm.HandleEvent(event)

		if event.Key() == tcell.KeyCtrlF {
			app.SetFocus(filterEditor)
			return event
		}

		if event.Key() == tcell.KeyCtrlS {
			app.SetFocus(sortEditor)
			return event
		}

		if event.Key() == tcell.KeyCtrlP {
			app.SetFocus(projectEditor)
			return event
		}

		if event.Key() == tcell.KeyCtrlU {
			return handleUseDatabaseEvent(event)
		}

		err := databaseTree.HandleDiconnectionEvent(event, func(key string) error {
			return disconnect(getConnection(), key)
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

	app.SetRoot(pages, true).EnableMouse(true)
}

// GetApplication returns the initialized *tview.Application.
func GetApplication() *tview.Application {
	return app
}

func buildMainScreen() {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandsView, 5, 1, false).
		AddItem(tview.NewFlex().
			AddItem(databaseTree, 40, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(editorView, 5, 1, false).
				AddItem(resultView, 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Statistics"), 3, 1, false),
				0, 1, false),
			0, 1, false)

	frame := tview.NewFrame(flex).
		AddText("MongoTUI - MongoDB crawler", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Copyright 2021-2024 Jan Renken", true, tview.AlignRight, tcell.ColorGreenYellow)
	pages.AddPage("frame", frame, true, true)
}

func setupCommandsView() {
	for i, command := range settings.GetCommands() {
		separator := "\t"
		if (i+1)%5 == 0 {
			separator = "\n"
		}
		fmt.Fprintf(commandsView, "%s%s", command.Description, separator)
	}
}
