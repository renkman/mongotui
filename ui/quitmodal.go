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

func GetQuitModalWidget(app *tview.Application, pages *tview.Pages) *ModalWidget {
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
