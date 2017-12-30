package hambidgetree

type Tree struct {
	uniqueId     NodeID
	ratios       TreeRatios
	xyRatioIndex int
	zyRatioIndex int
	scale        float64
	root         *Node
	epsilon      float64
}

func NewTree2D(ratios TreeRatios, ratioIndex int) *Tree {
	return NewTree(ratios, ratioIndex, -1)
}

func NewTree(ratios TreeRatios, xyRatioIndex int, zyRatioIndex int) *Tree {
	if xyRatioIndex < 0 || xyRatioIndex >= ratios.Ratios().Len() {
		panic("Invalid ratio index")
	}

	if zyRatioIndex >= ratios.Ratios().Len() {
		panic("Invalid ratio index")
	}

	tree := &Tree{
		uniqueId:     0,
		ratios:       ratios,
		xyRatioIndex: xyRatioIndex,
		zyRatioIndex: zyRatioIndex,
		scale:        1.0,
		epsilon:      CalculateRatiosEpsilon(ratios.Ratios()),
	}

	tree.root = NewNode(tree, nil)
	return tree
}

func (tree *Tree) nextNodeId() NodeID {
	id := tree.uniqueId
	tree.uniqueId = tree.uniqueId + 1
	return id
}

func (tree *Tree) RatioIndex(node *Node, plane RatioPlane) int {
	if node.parent == nil {
		if plane == RatioPlaneXY {
			return tree.xyRatioIndex
		} else {
			return tree.zyRatioIndex
		}
	}

	return node.RatioIndex()
}

func (tree *Tree) Ratio(ratioIndex int) float64 {
	if ratioIndex < 0 {
		return 0
	}

	return tree.ratios.Ratios().At(ratioIndex)
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

func (tree *Tree) RatioXY() float64 {
	return tree.ratios.Ratios().At(tree.xyRatioIndex)
}

func (tree *Tree) RatioZY() float64 {
	if tree.zyRatioIndex < 0 {
		return 0
	}

	return tree.ratios.Ratios().At(tree.zyRatioIndex)
}

// Useful for testing, each level splits the tree in half, starting with a
// square
// Even levels split all leaves in half using a vertical line
// Odd levels split those splits in half with a horizontal line creating squares again
func NewGridTree2D(levels int) *Tree {
	ratios := NewRatios([]float64{0.5, 1.0})
	treeRatios := NewTreeRatios(ratios, 0.0000001)

	tree := NewTree2D(treeRatios, 1)

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

func NewGridTree3D(levels int) *Tree {
	ratios := NewRatios([]float64{0.5, 1.0})
	treeRatios := NewTreeRatios(ratios, 0.0000001)

	tree := NewTree(treeRatios, 1, 1)

	for i := 0; i < levels; i++ {
		leaves := tree.Leaves()
		for _, leaf := range leaves {
			var split Split
			index := i % 3
			switch index {
			case 0:
				split = NewVerticalSplit(0, 0)
			case 1:
				split = NewHorizontalSplit(1, 1)
			case 2:
				split = NewDepthSplit(1, 1)
			}

			leaf.Divide(split)
		}
	}

	return tree
}
