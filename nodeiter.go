package hambidgetree

type NodeIterator struct {
	nodes []Node
}

func NewNodeIterator(root Node) *NodeIterator {
	return &NodeIterator{
		nodes: []Node{root},
	}
}

func (it *NodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

func (it *NodeIterator) Next() Node {
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

func FindLeaves(tree Tree) []Node {
	var leaves []Node
	it := NewNodeIterator(tree.Root())
	for it.HasNext() {
		node := it.Next()
		if node.Branch() == nil {
			leaves = append(leaves, node)
		}
	}
	return leaves
}
