package hambidgetree

type Tree struct {
	ratios     TreeRatios
	ratioIndex int
	scale      float64
	root       *Node
}

func NewTree(ratios TreeRatios, ratioIndex int) *Tree {
	if ratioIndex < 0 || ratioIndex >= ratios.Ratios().Len() {
		panic("Invalid ratio index")
	}

	tree := &Tree{
		ratios:     ratios,
		ratioIndex: ratioIndex,
		scale:      1.0,
	}

	tree.root = NewNode(tree, nil)
	return tree
}

func (tree *Tree) Leaves() []*Node {
	return tree.FilterNodes(func(node *Node) bool {
		return node.IsLeaf()
	})
}

func (tree *Tree) FilterNodes(filter func(*Node) bool) []*Node {
	var nodes []*Node
	it := NewNodeIterator(tree.root)
	for it.HasNext() {
		node := it.Next()
		if filter(node) {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree *Tree) Ratio() float64 {
	return tree.ratios.Ratios().At(tree.ratioIndex)
}
