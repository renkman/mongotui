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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResultTreeWidget struct {
	*tview.TreeView
	*Widget
}

func createResultTree() *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle("Result")

	return tree
}

func CreateResultTree(app *tview.Application,
	pages *tview.Pages) *ResultTreeWidget {
	tree := createResultTree()
	widget := createWidget(tree, "resulttree", app, pages)

	treeWidget := ResultTreeWidget{tree, widget}

	return &treeWidget
}

func (r *ResultTreeWidget) SetResult(result interface{}) {
	root := tview.NewTreeNode("Result Collection")
	r.TreeView.SetRoot(root).SetCurrentNode(root)

	addNode(root, result)
}

func addNode(node *tview.TreeNode, value interface{}) {
	switch value.(type) {
	case primitive.M:
		resultMap := value.(primitive.M)
		for k, v := range resultMap {
			child := tview.NewTreeNode(k)
			node.AddChild(child)
			addNode(child, v)
		}
	case primitive.D:
		resultMap := value.(primitive.D)
		for i, v := range resultMap {
			child := tview.NewTreeNode(fmt.Sprintf("%v", i))
			node.AddChild(child)
			addNode(child, v)
		}
	case primitive.E:
		resultElement := value.(primitive.E)
		child := tview.NewTreeNode(resultElement.Key)
		node.AddChild(child)
		addNode(child, resultElement.Value)
	default:
		node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", value)))
	}
}
