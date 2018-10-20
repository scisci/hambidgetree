package algo

import (
	htree "github.com/scisci/hambidgetree"
)

type NodeIterator struct {
	nodes []htree.Node
}

func NewNodeIterator(root htree.Node) *NodeIterator {
	return &NodeIterator{
		nodes: []htree.Node{root},
	}
}

func (it *NodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

func (it *NodeIterator) Next() htree.Node {
	if !it.HasNext() {
		return nil
	}

	node := it.nodes[len(it.nodes)-1]
	it.nodes = it.nodes[:len(it.nodes)-1]

	branch := node.Branch()
	if branch != nil {
		it.nodes = append(it.nodes, branch.Right(), branch.Left())
	}

	return node
}
