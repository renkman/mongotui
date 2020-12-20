package ui

import (
	// "github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type position struct{ x, y int }

type Textarea struct {
	textview   tview.TextView
	buffer     []byte
	cursorPos  position
	showCursor bool
}

func (t Textarea) append() {
	// if t.HasFocus() {
	// 	switch event.Key() {
	// 	case tcell.KeyRune:
	// 		t.textview.Write([]byte(string(event.Rune())))
	// 	case tcell.KeyEnter:
	// 		t.textview.Write([]byte("\n"))
	// 	}
	// 	return event
	// }
}

//editor := tview.NewTextView().SetWrap(true)

// editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 	if editor.HasFocus() {SwitchToPage
// 		switch event.Key() {
// 		case tcell.KeyRune:
// 			editor.Write([]byte(string(event.Rune())))
// 		case tcell.KeyEnter:
// 			editor.Write([]byte("\n"))
// 		}
// 	}
// 	return event
// })
