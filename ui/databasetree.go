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
	"github.com/rivo/tview"
)

// DatabaseTreeWidget displays the current MongoDB connections, including the
// databases and their collections in a tree view.
type DatabaseTreeWidget struct {
	*tview.TreeView
	*EventWidget
	loadCollections func(connectionURI string, database string) []string
}

const (
	rootText            string = "Connections"
	nodeLevelDatabase   string = "database"
	nodeLevelCollection string = "collection"
)

var parentMapping map[*tview.TreeNode]*tview.TreeNode = make(map[*tview.TreeNode]*tview.TreeNode)

// CreateDatabaseTreeWidget creates a new DatabaseTreeWidget.
func CreateDatabaseTreeWidget(app *tview.Application,
	pages *tview.Pages,
	loadCollections func(connectionURI string, database string) []string) *DatabaseTreeWidget {
	tree := createDatabaseTree()
	widget := createEventWidget(tree, "databasetree", tcell.KeyCtrlD, app, pages)

	treeWidget := DatabaseTreeWidget{tree, widget, loadCollections}
	tree.SetSelectedFunc(treeWidget.addCollections)

	return &treeWidget
}

// SetFocus implements the FocusSetter interface to set the focus to the
// tview.TreeView.
func (d *DatabaseTreeWidget) SetFocus(app *tview.Application) {
	app.SetFocus(d)
}

// HandleEvent sets the event key of the DatabaseTreeWidget.
func (d *DatabaseTreeWidget) HandleEvent(event *tcell.EventKey) {
	d.setEvent(d, event, false)
}

// HandleDiconnectionEvent disconnects from the selected instance if a client node is
// selected.
func (d *DatabaseTreeWidget) HandleDiconnectionEvent(event *tcell.EventKey, disconnect func(clientKey string) error) error {
	if event.Key() != tcell.KeyCtrlT {
		return nil
	}
	node := d.GetCurrentNode()
	reference := node.GetReference()
	if reference == nil ||
		reference.(string) == nodeLevelDatabase ||
		reference.(string) == nodeLevelCollection {
		return nil
	}
	err := disconnect(reference.(string))
	if err != nil {
		return err
	}

	root := d.GetRoot()
	connections := root.GetChildren()
	root.ClearChildren()
	for _, connection := range connections {
		if connection != node {
			root.AddChild(connection)
		}
	}
	return nil
}

// AddDatabases adds the databases of the instance of the passed connectionURI to the
// connection tree.
func (d *DatabaseTreeWidget) AddDatabases(host string, connectionURI string, databases []string) {
	root := d.TreeView.GetRoot()
	var connectionNode *tview.TreeNode
	for _, node := range root.GetChildren() {
		if node.GetReference().(string) == connectionURI {
			connectionNode = node
			break
		}
	}
	if connectionNode == nil {
		if host == "" {
			host = connectionURI
		}
		connectionNode = tview.NewTreeNode(host).
			SetColor(tcell.ColorGreenYellow).
			SetReference(connectionURI)
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

// UpdateCollections removes and re-adds the collections of the selected database.
func (d *DatabaseTreeWidget) UpdateCollections() {
	currentNode := d.TreeView.GetCurrentNode()
	if currentNode == nil {
		return
	}
	d.addCollections(currentNode)
}

func createDatabaseTree() *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("Databases")

	root := tview.NewTreeNode(rootText)
	tree.SetRoot(root).SetCurrentNode(root)

	return tree
}

func (d *DatabaseTreeWidget) getCollections(connectionURI string, name string) []string {
	return d.loadCollections(connectionURI, name)
}

func (d *DatabaseTreeWidget) addCollections(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil || reference.(string) != nodeLevelDatabase {
		return
	}

	node.ClearChildren()
	connectionURI := parentMapping[node].GetReference().(string)
	collections := d.getCollections(connectionURI, node.GetText())

	for _, collection := range collections {
		node.AddChild(tview.NewTreeNode(collection).
			SetReference(nodeLevelCollection))
	}
}
