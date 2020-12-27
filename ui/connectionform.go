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

func createConnectionForm(cancel func()) (*tview.Flex, *tview.Form) {
	connectionForm := tview.NewForm().
		AddInputField("Host:", "", 20, nil, nil).
		AddInputField("Port:", "", 20, tview.InputFieldInteger, nil).
		AddInputField("User:", "", 20, nil, nil).
		AddPasswordField("Password:", "", 20, '*', nil).
		AddInputField("Replicaset:", "", 20, nil, nil).
		AddCheckbox("TLS/SSL:", false, nil).
		AddButton("Connect", nil).
		AddButton("Cancel", cancel)
	connectionForm.SetBorder(true).SetTitle("Mongo DB Connection")

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(connectionForm, 18, 1, false).
			AddItem(nil, 0, 1, false), 40, 1, false).
		AddItem(nil, 0, 1, false)

	return modal, connectionForm
}

func GetConnectionFormWidget(app *tview.Application, pages *tview.Pages, connect func(connection *models.Connection)) *FormWidget {
	modal, form := createConnectionForm(
		func() {
			pages.RemovePage("connection")
		})
	widget := createEventWidget(modal, "connection", tcell.KeyCtrlC, app, pages)
	formWidget := FormWidget{modal, form, widget}

	formWidget.SetSelectedFunc(connect)

	return &formWidget
}

func (f *FormWidget) SetFocus(app *tview.Application) {
	app.SetFocus(f.GetFormItem(0))
}

func (f *FormWidget) SetEvent(event *tcell.EventKey) {
	f.setEvent(f, event)
}

func (f *FormWidget) SetSelectedFunc(connect func(connection *models.Connection)) {
	f.GetButton(0).SetSelectedFunc(func() {
		connection := f.GetData()
		connect(&connection)
	})
}

func (f *FormWidget) GetData() models.Connection {
	connection := models.Connection{}
	connection.Host = f.GetFormItem(0).(*tview.InputField).GetText()
	connection.Port = f.GetFormItem(1).(*tview.InputField).GetText()
	connection.User = f.GetFormItem(2).(*tview.InputField).GetText()
	connection.Password = f.GetFormItem(3).(*tview.InputField).GetText()
	connection.Replicaset = f.GetFormItem(4).(*tview.InputField).GetText()

	return connection
}
