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
	"github.com/renkman/mongotui/models"
	"github.com/rivo/tview"
)

type FormWidget struct {
	*tview.Flex
	*tview.Form
	*EventWidget
}

func createConnectionForm(cancel func(),
	getSavedConnections func() ([]string, error),
	loadSavedConnection func(string) (string, error)) (*tview.Flex, *tview.Form) {
	savedConnections, err := getSavedConnections()
	if err != nil {
		savedConnections = []string{}
	}

	connectionForm := tview.NewForm().
		AddInputField("Host:", "", 20, nil, nil).
		AddInputField("Port:", "", 20, tview.InputFieldInteger, nil).
		AddInputField("User:", "", 20, nil, nil).
		AddPasswordField("Password:", "", 20, '*', nil).
		AddInputField("Replicaset:", "", 20, nil, nil).
		AddCheckbox("TLS/SSL:", false, nil).
		AddInputField("URI:", "", 40, nil, nil).
		AddCheckbox("Save connection:", false, nil).
		AddButton("Connect", nil).
		AddButton("Cancel", cancel)
	connectionForm.SetBorder(true).SetTitle("Mongo DB Connection")

	connectionDropDown := tview.NewDropDown().
		SetLabel("Load saved connection: ").
		SetOptions(savedConnections, nil).
		SetSelectedFunc(func(key string, index int) {
			connectionURI, _ := loadSavedConnection(key)
			connectionForm.GetFormItem(6).(*tview.InputField).SetText(connectionURI)
		})

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(connectionDropDown, 2, 1, false).
		AddItem(connectionForm, 23, 1, false)

	frame := tview.NewFrame(formFlex).
		AddText("Set fields individually or directly set the URI.", true, tview.AlignCenter, tcell.ColorYellow).
		AddText("If the URI is set, the individual fields are ignored.", true, tview.AlignCenter, tcell.ColorYellow)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(frame, 30, 1, false).
			AddItem(nil, 0, 1, false), 60, 1, false).
		AddItem(nil, 0, 1, false)

	return modal, connectionForm
}

func CreateConnectionFormWidget(app *tview.Application, pages *tview.Pages, connect func(connection *models.Connection), getSavedConnections func() ([]string, error), loadSavedConnection func(string) (string, error)) *FormWidget {
	modal, form := createConnectionForm(
		func() {
			pages.RemovePage("connection")
		},
		getSavedConnections,
		loadSavedConnection)
	widget := createEventWidget(modal, "connection", tcell.KeyCtrlC, app, pages)
	formWidget := FormWidget{modal, form, widget}

	formWidget.setSelectedFunc(connect)

	return &formWidget
}

func (f *FormWidget) SetFocus(app *tview.Application) {
	app.SetFocus(f.GetFormItem(0))
}

func (f *FormWidget) SetEvent(event *tcell.EventKey) {
	f.setEvent(f, event)
}

func (f *FormWidget) setSelectedFunc(connect func(connection *models.Connection)) {
	f.GetButton(0).SetSelectedFunc(func() {
		connection := f.getData()
		connect(&connection)
	})
}

func (f *FormWidget) getData() models.Connection {
	connection := models.Connection{}
	connection.Host = f.GetFormItem(0).(*tview.InputField).GetText()
	connection.Port = f.GetFormItem(1).(*tview.InputField).GetText()
	connection.User = f.GetFormItem(2).(*tview.InputField).GetText()
	connection.Password = f.GetFormItem(3).(*tview.InputField).GetText()
	connection.Replicaset = f.GetFormItem(4).(*tview.InputField).GetText()
	connection.TLS = f.GetFormItem(5).(*tview.Checkbox).IsChecked()
	connection.URI = f.GetFormItem(6).(*tview.InputField).GetText()
	connection.SaveConnection = f.GetFormItem(7).(*tview.Checkbox).IsChecked()

	return connection
}
