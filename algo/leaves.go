package algo

import (
	htree "github.com/scisci/hambidgetree"
)

func FindLeaves(tree htree.Tree) []htree.Node {
	var leaves []htree.Node
	it := htree.NewNodeIterator(tree.Root())
	for it.HasNext() {
		node := it.Next()
		if node.Branch() == nil {
			leaves = append(leaves, node)
		}
	}
	return leaves
}
