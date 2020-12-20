package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/renkman/mongotui/models"
	"github.com/renkman/mongotui/mongo"
	"github.com/renkman/mongotui/ui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	textView := tview.NewTextView()
	editor := tview.NewInputField().
		SetLabel("Command: ").
		SetFieldWidth(200)
	editor.SetDoneFunc(func(key tcell.Key) {
		textView.Write([]byte(editor.GetText()))
	})

	commandsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter)

	commandsView.SetBorder(true).
		SetTitle("Commands")

	fmt.Fprint(commandsView, "\n[white]Ctrl - Q [darkcyan]Quit\t[white]Ctrl - C[darkcyan] Connect to database")

	textView.SetBorder(true).SetTitle("Result")
	editor.SetBorder(true).SetTitle("Editor")

	ctx := context.Background()
	databaseTree := ui.GetDatabaseTree(app, pages, func(name string) []string {
		mongo.UseDatabase(name)
		collections, err := mongo.GetCollections(ctx)
		if err != nil {
			message := fmt.Sprintf("Getting collections of database %s failed:\n\n%s", name, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return collections
		}
		return collections
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandsView, 5, 1, false).
		AddItem(tview.NewFlex().
			AddItem(databaseTree, 40, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(editor, 3, 1, false).
				AddItem(textView, 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Statistics"), 10, 1, false),
				0, 1, false),
			0, 1, false)

	frame := tview.NewFrame(flex).AddText("MongoTUI - MongoDB crawler", true, tview.AlignLeft, tcell.ColorYellow)
	pages.AddPage("frame", frame, true, true)

	quitModal := ui.GetQuitModalWidget(app, pages)

	connectionForm := ui.GetConnectionFormWidget(app, pages, func(connection models.Connection) {
		err := mongo.Connect(ctx, connection)
		if err != nil {
			message := fmt.Sprintf("Connection to %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}
		//defer mongo.Disconnect(ctx)

		databases, err := mongo.GetDatabases(ctx)
		if err != nil {
			message := fmt.Sprintf("Getting databeses of %s failed:\n\n%s", connection.Host, err.Error())
			ui.CreateMessageModalWidget(app, pages, ui.TypeError, message)
			return
		}
		databaseTree.AddDatabases(connection.Host, databases)
		pages.RemovePage("connection")
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quitModal.SetEvent(event)
		connectionForm.SetEvent(event)
		databaseTree.SetEvent(event)

		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyNUL, ' ', tcell.ModNone)
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).SetFocus(editor).Run(); err != nil {
		e := fmt.Sprint(err)
		fmt.Print(e)
		panic(err)
	}
}
