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
	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/models"
	"github.com/rivo/tview"
)

const connectionFormWidget string = "connection"

// FormWidget is a MongoDB connection data form modal. It contains individual fields,
// a mongodb:// connection URI field to enter connection data, a checkbox for storing
// the connection data locally and a dropbox to select a former stored connection.
type FormWidget struct {
	*tview.Flex
	*tview.Form
	*tview.DropDown
	*EventWidget
	getSavedConnections func() ([]string, error)
	loadSavedConnection func(string) (string, error)
}

// CreateConnectionFormWidget creates a new FormWidget.
func CreateConnectionFormWidget(app *tview.Application, pages *tview.Pages, connect func(connection *models.Connection), canStoreConnection bool, getSavedConnections func() ([]string, error), loadSavedConnection func(string) (string, error)) *FormWidget {
	modal, form, dropDown := createConnectionForm(
		func() {
			pages.RemovePage(connectionFormWidget)
		},
		canStoreConnection)
	widget := createEventWidget(modal, connectionFormWidget, tcell.KeyCtrlC, app, pages)
	formWidget := FormWidget{modal, form, dropDown, widget, getSavedConnections, loadSavedConnection}

	formWidget.setButtonSelectedFunc(connect)

	return &formWidget
}

// SetFocus implements the FocusSetter interface to set the focus tview.Form when
// called.
func (f *FormWidget) SetFocus(app *tview.Application) {
	savedConnections, err := f.getSavedConnections()
	if err != nil {
		savedConnections = []string{}
	}
	f.SetOptions(savedConnections, f.setDropBoxSelectedFunc)

	app.SetFocus(f.Form)
}

// HandleEvent handles the event key of the FormWidget.
func (f *FormWidget) HandleEvent(event *tcell.EventKey) {
	f.handleEvent(f, event, true)
}

func createConnectionForm(cancel func(),
	canStoreConnection bool) (*tview.Flex, *tview.Form, *tview.DropDown) {

	connectionForm := tview.NewForm().
		AddInputField("Host:", "", 20, nil, nil).
		AddInputField("Port:", "", 20, tview.InputFieldInteger, nil).
		AddInputField("User:", "", 20, nil, nil).
		AddPasswordField("Password:", "", 20, '*', nil).
		AddInputField("Replicaset:", "", 20, nil, nil).
		AddCheckbox("TLS/SSL:", false, nil).
		AddInputField("URI:", "", 40, nil, nil)
	if canStoreConnection {
		connectionForm.AddCheckbox("Save connection:", false, nil)
	}
	connectionForm.AddButton("Connect", nil).
		AddButton("Cancel", cancel)
	connectionForm.SetBorder(true).SetTitle("Mongo DB connection")

	connectionDropDown := tview.NewDropDown().
		SetLabel("Load saved connection: ")

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	if canStoreConnection {
		formFlex.AddItem(connectionDropDown, 2, 1, false)
	}
	formFlex.AddItem(connectionForm, 23, 1, false)

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

	return modal, connectionForm, connectionDropDown
}

func (f *FormWidget) setButtonSelectedFunc(connect func(connection *models.Connection)) {
	f.GetButton(0).SetSelectedFunc(func() {
		connection := f.getData()
		f.Pages.RemovePage(connectionFormWidget)
		go connect(&connection)
	})
}

func (f *FormWidget) setDropBoxSelectedFunc(key string, index int) {
	connectionURI, _ := f.loadSavedConnection(key)
	f.GetFormItem(6).(*tview.InputField).SetText(connectionURI)
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
	if f.GetFormItemCount() == 8 {
		connection.SaveConnection = f.GetFormItem(7).(*tview.Checkbox).IsChecked()
	}

	return connection
}
