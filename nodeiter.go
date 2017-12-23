package hambidgetree

type NodeIterator struct {
	nodes []*Node
}

func NewNodeIterator(root *Node) *NodeIterator {
	return &NodeIterator{
		nodes: []*Node{root},
	}
}

func (it *NodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

func (it *NodeIterator) Next() *Node {
	if !it.HasNext() {
		return nil
	}

	node := it.nodes[len(it.nodes)-1]
	it.nodes = it.nodes[:len(it.nodes)-1]

	if !node.IsLeaf() {
		it.nodes = append(it.nodes, node.right, node.left)
	}

	return node
}
