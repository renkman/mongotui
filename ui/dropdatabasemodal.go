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

package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const dropdatabasemodal string = "dropdatabase"

// DropDatabaseModalWidget is the modal to drop a database.
type DropDatabaseModalWidget struct {
	*tview.Modal
	*EventWidget
	getCurrentDatabase func() string
}

// CreateDropDatabaseModalWidget creates a new DropDatabaseModalWidget.
func CreateDropDatabaseModalWidget(app *tview.Application, pages *tview.Pages,
	getCurrentDatabase func() string, drop func()) *DropDatabaseModalWidget {
	modal := createDropDatabaseModal(
		func() {
			drop()
			pages.RemovePage(dropdatabasemodal)
		},
		func() {
			pages.RemovePage(dropdatabasemodal)
		})
	widget := createEventWidget(modal, dropdatabasemodal, tcell.KeyCtrlX, app, pages)
	return &DropDatabaseModalWidget{modal, widget, getCurrentDatabase}
}

// SetFocus implements the FocusSetter interface to set the focus to the OK button.
func (m *DropDatabaseModalWidget) SetFocus(app *tview.Application) {
	name := m.getCurrentDatabase()
	m.SetText(fmt.Sprintf("Do you really want to drop %s?", name))
	app.SetFocus(m)
}

// HandleEvent handles the event key of the DropDatabaseModalWidget.
func (m *DropDatabaseModalWidget) HandleEvent(event *tcell.EventKey) {
	m.handleEvent(m, event, true)
}

func createDropDatabaseModal(drop func(), cancel func()) *tview.Modal {
	quitModal := tview.NewModal().
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				drop()
			} else {
				cancel()
			}
		})
	return quitModal
}
