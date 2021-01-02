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

	"github.com/rivo/tview"
)

const name = "messageModal"

// MessageType indicates the type of a message.
type MessageType int

// Currently supported message types.
const (
	TypeInfo  MessageType = 0
	TypeError MessageType = 1
)

var messageTypes = [2]string{"Info", "Error"}

// MessageModalWidget is a simple message dialog, which displays the message type,
// a text message and an Ok button.
type MessageModalWidget struct {
	*tview.Modal
	*Widget
	MessageType MessageType
	Message     string
	Name        string
}

type unknownMessageTypeError struct {
	InvalidValue MessageType
}

func (e *unknownMessageTypeError) Error() string {
	return fmt.Sprintf("Invalid message type %v", e.InvalidValue)
}

// CreateMessageModalWidget creates a new MessageModalWidget.
func CreateMessageModalWidget(app *tview.Application, pages *tview.Pages, messageType MessageType, message string) (*MessageModalWidget, error) {
	if len(messageTypes) <= int(messageType) || len(messageTypes) < 0 {
		return nil, &unknownMessageTypeError{messageType}
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

func createMessageModal(messageType MessageType, message string, ok func()) *tview.Modal {
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
