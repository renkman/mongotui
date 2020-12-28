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

type ModalWidget struct {
	*tview.Modal
	*EventWidget
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

func CreateQuitModalWidget(app *tview.Application, pages *tview.Pages) *ModalWidget {
	modal := createQuitModal(
		func() { app.Stop() },
		func() {
			pages.RemovePage("quit")
		})
	widget := createEventWidget(modal, "quit", tcell.KeyCtrlQ, app, pages)
	return &ModalWidget{modal, widget}
}

func (m *ModalWidget) SetFocus(app *tview.Application) {
	app.SetFocus(m)
}

func (m *ModalWidget) SetEvent(event *tcell.EventKey) {
	m.setEvent(m, event)
}
