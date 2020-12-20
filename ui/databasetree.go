package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type DatabaseTreeWidget struct {
	*tview.TreeView
	*EventWidget
	loadCollections func(database string) []string
}

const rootText string = "Connections"
const nodeLevelDatabase string = "database"
const nodeLevelCollection string = "collection"

func createDatabaseTree() *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("Databases")

	root := tview.NewTreeNode(rootText)
	tree.SetRoot(root).SetCurrentNode(root)

	return tree
}

func GetDatabaseTree(app *tview.Application,
	pages *tview.Pages,
	loadCollections func(database string) []string) *DatabaseTreeWidget {
	tree := createDatabaseTree()
	widget := createEventWidget(tree, "databasetree", tcell.KeyCtrlD, app, pages)

	treeWidget := DatabaseTreeWidget{tree, widget, loadCollections}
	tree.SetSelectedFunc(treeWidget.addCollections)

	return &treeWidget
}

func (d *DatabaseTreeWidget) SetFocus(app *tview.Application) {
	app.SetFocus(d)
}

func (d *DatabaseTreeWidget) SetEvent(event *tcell.EventKey) {
	d.setEvent(d, event)
}

func (d *DatabaseTreeWidget) AddDatabases(connection string, databases []string) {
	connectionNode := tview.NewTreeNode(connection).SetColor(tcell.ColorGreen)
	d.TreeView.GetRoot().AddChild(connectionNode)

	for _, database := range databases {
		connectionNode.AddChild(tview.NewTreeNode(database).
			SetColor(tcell.ColorYellow).
			SetSelectable(true).
			SetReference(nodeLevelDatabase))
	}
}

func (d *DatabaseTreeWidget) getCollections(name string) []string {
	return d.loadCollections(name)
}

func (d *DatabaseTreeWidget) addCollections(node *tview.TreeNode) {
	reference := node.GetReference()
	if node.GetText() == rootText || reference == nil {
		return
	}

	if reference.(string) != nodeLevelDatabase {
		return
	}

	node.ClearChildren()
	collections := d.getCollections(node.GetText())

	for _, collection := range collections {
		node.AddChild(tview.NewTreeNode(collection)).
			SetReference(nodeLevelCollection)
	}
}
