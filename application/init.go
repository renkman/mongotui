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
	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/settings"
	"github.com/renkman/mongotui/ui"
	"github.com/rivo/tview"
)

var (
	app            *tview.Application
	pages          *tview.Pages
	connectionForm *ui.FormWidget
	databaseTree   *ui.DatabaseTreeWidget
)

func init() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	databaseTree = ui.CreateDatabaseTreeWidget(app, pages, updateDatabaseTree)

	connectionForm := ui.CreateConnectionFormWidget(app, pages, connect, settings.CanStoreConnection, settings.GetConnections, settings.GetConnectionURI)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		connectionForm.HandleEvent(event)
		databaseTree.HandleEvent(event)
		return event
	})

	app.SetRoot(pages, true).EnableMouse(true)
}

func getApplication() *tview.Application {
	return app
}
