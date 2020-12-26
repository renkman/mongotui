package ui

import (
	"fmt"
	//"reflect"

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
	case primitive.D:
		resultMap := value.(primitive.D).Map()
		for k, v := range resultMap {
			child := tview.NewTreeNode(k)
			node.AddChild(child)
			addNode(child, v)
		}
	default:
		node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", value)))
	}
}
