package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type FocusSetter interface {
	SetFocus(app *tview.Application)
}

type Widget struct {
	Primitive tview.Primitive
	Name      string
	App       *tview.Application
	Pages     *tview.Pages
}

type EventWidget struct {
	*Widget
	Key tcell.Key
}

func createWidget(primitive tview.Primitive, name string, app *tview.Application, pages *tview.Pages) *Widget {
	w := Widget{primitive, name, app, pages}
	return &w
}

func createEventWidget(primitive tview.Primitive, name string, key tcell.Key, app *tview.Application, pages *tview.Pages) *EventWidget {
	w := EventWidget{&Widget{primitive, name, app, pages}, key}
	return &w
}

func (w *EventWidget) setEvent(f FocusSetter, event *tcell.EventKey) {
	if event.Key() != w.Key {
		return
	}
	w.Pages.AddPage(w.Name, w.Primitive, true, true)
	f.SetFocus(w.App)
}
