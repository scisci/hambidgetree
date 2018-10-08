package hambidgetree

type ImmutableNodeIterator struct {
	nodes []ImmutableNode
}

func NewImmutableNodeIterator(root ImmutableNode) *ImmutableNodeIterator {
	return &ImmutableNodeIterator{
		nodes: []ImmutableNode{root},
	}
}

func (it *ImmutableNodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

func (it *ImmutableNodeIterator) Next() ImmutableNode {
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

func FindLeaves(tree ImmutableTree) []ImmutableNode {
	var leaves []ImmutableNode
	it := NewImmutableNodeIterator(tree.Root())
	for it.HasNext() {
		node := it.Next()
		if node.Branch() == nil {
			leaves = append(leaves, node)
		}
	}
	return leaves
}
