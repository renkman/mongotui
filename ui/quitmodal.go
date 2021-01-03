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
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// QuitModalWidget is the modal to quit the application. It displays the question
// whether the user really wants to quit and an Ok and a Cancel button.
type QuitModalWidget struct {
	*tview.Modal
	*EventWidget
}

// CreateQuitModalWidget creates a new QuitModalWidget.
func CreateQuitModalWidget(app *tview.Application, pages *tview.Pages) *QuitModalWidget {
	modal := createQuitModal(
		func() { app.Stop() },
		func() {
			pages.RemovePage("quit")
		})
	widget := createEventWidget(modal, "quit", tcell.KeyCtrlQ, app, pages)
	return &QuitModalWidget{modal, widget}
}

// SetFocus implements the FocusSetter interface to set the focus to the OK button.
func (m *QuitModalWidget) SetFocus(app *tview.Application) {
	app.SetFocus(m)
}

// SetEvent sets the event key of the QuitModalWidget.
func (m *QuitModalWidget) SetEvent(event *tcell.EventKey) {
	m.setEvent(m, event, true)
}

func createQuitModal(quit func(), cancel func()) *tview.Modal {
	quitModal := tview.NewModal().
		SetText("Do you really want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				quit()
			} else {
				cancel()
			}
		})
	return quitModal
}
