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
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResultTreeWidget is the tree view to display the MondoDB command results.
type ResultTreeWidget struct {
	*tview.TreeView
	*EventWidget
	loadNextResults     func()
	loadPreviousResults func()
}

// CreateResultTreeWidget creates a new ResultTreeWidget.
func CreateResultTreeWidget(app *tview.Application, pages *tview.Pages, loadNextResults func(), loadPreviousResults func()) *ResultTreeWidget {
	tree := createResultTree()
	widget := createEventWidget(tree, "resulttree", tcell.KeyCtrlR, app, pages)

	treeWidget := ResultTreeWidget{tree, widget, loadNextResults, loadPreviousResults}

	return &treeWidget
}

// SetResult sets the result returned by the executed MongoDB command to the
// tcell.TreeView.
func (r *ResultTreeWidget) SetResult(result []map[string]interface{}) {
	root := tview.NewTreeNode("Result Collection")
	r.TreeView.SetRoot(root).SetCurrentNode(root)

	for _, document := range result {
		addNode(root, document)
	}
}

// SetFocus implements the FocusSetter interface to set the focus to the
// tview.TreeView.
func (r *ResultTreeWidget) SetFocus(app *tview.Application) {
	app.SetFocus(r)
}

// HandleEvent handles the event key of the ResultTreeWidget.
func (r *ResultTreeWidget) HandleEvent(event *tcell.EventKey, app *tview.Application) {
	if app.GetFocus() == r && event.Key() == tcell.KeyRune {
		if event.Rune() == 'n' {
			r.loadNextResults()
			return
		}
		if event.Rune() == 'p' {
			r.loadPreviousResults()
			return
		}
	}

	r.handleEvent(r, event, false)
}

func createResultTree() *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("Result")

	return tree
}

func addMapNode(node *tview.TreeNode, document map[string]interface{}) {
	for key, value := range document {
		child := tview.NewTreeNode(key)
		node.AddChild(child)
		addNode(child, value)
	}
}

func addNode(node *tview.TreeNode, value interface{}) {
	switch value.(type) {
	case map[string]interface{}:
		resultMap := value.(map[string]interface{})
		addMapNode(node, resultMap)
	case primitive.A:
		resultArray := value.(primitive.A)
		for i, v := range resultArray {
			child := tview.NewTreeNode(fmt.Sprintf("%v", i))
			node.AddChild(child)
			addNode(child, v)
		}
	default:
		node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", value)))
	}
}

// func addNode(node *tview.TreeNode, value interface{}) {
// 	switch value.(type) {
// 	case primitive.A:
// 		resultMap := value.(primitive.A)
// 		for i, v := range resultMap {
// 			child := tview.NewTreeNode(fmt.Sprintf("%v", i))
// 			node.AddChild(child)
// 			addNode(child, v)
// 		}
// 	case primitive.M:
// 		resultMap := value.(primitive.M)
// 		for k, v := range resultMap {
// 			child := tview.NewTreeNode(k)
// 			node.AddChild(child)
// 			addNode(child, v)
// 		}
// 	case primitive.D:
// 		resultMap := value.(primitive.D)
// 		for i, v := range resultMap {
// 			child := tview.NewTreeNode(fmt.Sprintf("%v", i))
// 			node.AddChild(child)
// 			addNode(child, v)
// 		}
// 	case primitive.E:
// 		resultElement := value.(primitive.E)
// 		child := tview.NewTreeNode(resultElement.Key)
// 		node.AddChild(child)
// 		addNode(child, resultElement.Value)
// 	default:
// 		node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", value)))
// 	}
// }
