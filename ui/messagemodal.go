package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

const name = "messageModal"
const TypeInfo = 0
const TypeError = 1

var messageTypes = [2]string{"Info", "Error"}

type MessageModalWidget struct {
	*tview.Modal
	*Widget
	MessageType int
	Message     string
	Name        string
}

type UnknownMessageTypeError struct {
	InvalidValue int
}

func (e *UnknownMessageTypeError) Error() string {
	return fmt.Sprintf("Invalid message type %v", e.InvalidValue)
}

func createMessageModal(messageType int, message string, ok func()) *tview.Modal {
	modalTypeText := messageTypes[messageType]
	messageText := fmt.Sprintf("%s\n\n%s", modalTypeText, message)

	modal := tview.NewModal().
		AddButtons([]string{"Ok"}).
		SetText(messageText).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ok()
		})

	return modal
}

func CreateMessageModalWidget(app *tview.Application, pages *tview.Pages, messageType int, message string) (*MessageModalWidget, error) {
	if len(messageTypes) <= messageType || len(messageTypes) < 0 {
		return nil, &UnknownMessageTypeError{messageType}
	}

	modal := createMessageModal(
		messageType,
		message,
		func() {
			pages.RemovePage(name)
		})
	widget := createWidget(modal, name, app, pages)

	pages.AddPage(name, modal, true, true)
	app.SetFocus(modal)

	return &MessageModalWidget{modal, widget, messageType, message, name}, nil
}
