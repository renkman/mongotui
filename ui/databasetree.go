package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type DatabaseTreeWidget struct {
	*tview.TreeView
	*EventWidget
	loadCollections func(connectionUri string, database string) []string
}

const rootText string = "Connections"
const nodeLevelDatabase string = "database"
const nodeLevelCollection string = "collection"

var parentMapping map[*tview.TreeNode]*tview.TreeNode = make(map[*tview.TreeNode]*tview.TreeNode)

func createDatabaseTree() *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("Databases")

	root := tview.NewTreeNode(rootText)
	tree.SetRoot(root).SetCurrentNode(root)

	return tree
}

func GetDatabaseTree(app *tview.Application,
	pages *tview.Pages,
	loadCollections func(connectionUri string, database string) []string) *DatabaseTreeWidget {
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

func (d *DatabaseTreeWidget) AddDatabases(host string, connectionUri string, databases []string) {
	root := d.TreeView.GetRoot()
	var connectionNode *tview.TreeNode
	for _, node := range root.GetChildren() {
		if node.GetReference().(string) == connectionUri {
			connectionNode = node
			break
		}
	}
	if connectionNode == nil {
		connectionNode = tview.NewTreeNode(host).
			SetColor(tcell.ColorGreen).
			SetReference(connectionUri)
		d.TreeView.GetRoot().AddChild(connectionNode)
	}

	connectionNode.ClearChildren()
	for _, database := range databases {
		databaseNode := tview.NewTreeNode(database).
			SetColor(tcell.ColorYellow).
			SetSelectable(true).
			SetReference(nodeLevelDatabase)
		connectionNode.AddChild(databaseNode)
		parentMapping[databaseNode] = connectionNode
	}
}

func (d *DatabaseTreeWidget) getCollections(connectionUri string, name string) []string {
	return d.loadCollections(connectionUri, name)
}

func (d *DatabaseTreeWidget) addCollections(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return
	}

	if reference.(string) != nodeLevelDatabase {
		return
	}

	node.ClearChildren()
	connectionUri := parentMapping[node].GetReference().(string)
	collections := d.getCollections(connectionUri, node.GetText())

	for _, collection := range collections {
		node.AddChild(tview.NewTreeNode(collection)).
			SetReference(nodeLevelCollection)
	}
}
