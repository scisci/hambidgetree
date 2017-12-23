package hambidgetree

type Tree struct {
	uniqueId   NodeID
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
		uniqueId:   0,
		ratios:     ratios,
		ratioIndex: ratioIndex,
		scale:      1.0,
	}

	tree.root = NewNode(tree, nil)
	return tree
}

func (tree *Tree) nextNodeId() NodeID {
	id := tree.uniqueId
	tree.uniqueId = tree.uniqueId + 1
	return id
}

func (tree *Tree) Leaves() []*Node {
	return tree.FilterNodes(func(node *Node) bool {
		return node.IsLeaf()
	})
}

func (tree *Tree) Root() *Node {
	return tree.root
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

// Useful for testing, each level splits the tree in half, starting with a
// square
// Even levels split all leaves in half using a vertical line
// Odd levels split those splits in half with a horizontal line creating squares again
func NewGridTree(levels int) *Tree {
	ratios := NewRatios([]float64{0.5, 1.0})
	treeRatios := NewTreeRatios(ratios, 0.0000001)

	tree := NewTree(treeRatios, 1)

	for i := 0; i < levels; i++ {
		leaves := tree.Leaves()
		for _, leaf := range leaves {
			var split Split
			if i&1 == 0 {
				split = NewVerticalSplit(0, 0)
			} else {
				split = NewHorizontalSplit(1, 1)
			}

			leaf.Divide(split)
		}
	}

	return tree
}
