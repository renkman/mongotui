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
	"github.com/rivo/tview"
)

// DatabaseFormWidget is a MongoDB connection data form modal. It contains individual fields,
// a mongodb:// connection URI field to enter connection data, a checkbox for storing
// the connection data locally and a dropbox to select a former stored connection.
type DatabaseFormWidget struct {
	*tview.Flex
	*tview.Form
	*Widget
}

const databaseFormWidget string = "dadatabaseForm"

// CreateDatabaseFormWidget creates a new DatabaseFormWidget.
func CreateDatabaseFormWidget(app *tview.Application, pages *tview.Pages, use func(name string)) *DatabaseFormWidget {
	modal, form := createDatabaseFormWidget(
		func() {
			pages.RemovePage(databaseFormWidget)
		})
	widget := createWidget(modal, databaseFormWidget, app, pages)
	formWidget := DatabaseFormWidget{modal, form, widget}

	formWidget.setButtonSelectedFunc(use)

	formWidget.Pages.AddPage(databaseFormWidget, modal, true, true)

	return &formWidget
}

// SetFocus implements the FocusSetter interface to set the focus tview.Form when
// called.
func (d *DatabaseFormWidget) SetFocus(app *tview.Application) {
	app.SetFocus(d.Form)
}

func createDatabaseFormWidget(cancel func()) (*tview.Flex, *tview.Form) {
	connectionForm := tview.NewForm().
		AddInputField("Name:", "", 20, nil, nil).
		AddButton("Use", nil).
		AddButton("Cancel", cancel)
	connectionForm.SetBorder(true).SetTitle("Use database")

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(connectionForm, 7, 1, false).
			AddItem(nil, 0, 1, false), 60, 1, false).
		AddItem(nil, 0, 1, false)

	return modal, connectionForm
}

func (d *DatabaseFormWidget) setButtonSelectedFunc(use func(name string)) {
	d.GetButton(0).SetSelectedFunc(func() {
		name := d.getName()
		use(name)
		d.Pages.RemovePage(databaseFormWidget)
	})
}

func (d *DatabaseFormWidget) getName() string {
	return d.GetFormItem(0).(*tview.InputField).GetText()
}
