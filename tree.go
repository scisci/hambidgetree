package hambidgetree

// Serialize this like so
//
// ratios: (tree.ratios.ratios) // should marshal themselves with some kind of type flag, if strings, then expressions, otherwise floats
// xyRatioIndex:
// zyRatioIndex:
// root: id
// node: [id: id, split: hvd, index, left: id, right: id]

/*
type ImmutableTree interface {
	Ratios() Ratios
	RatioIndexXY() int
	RatioIndexZY() int
	Root() ImmutableNode
}

type ImmutableNode interface {
	ID() NodeID
	SplitType() SplitType
	RatioIndex() int
	Left() ImmutableNode
	Right() ImmutableNode
}

type ImmutableBranch interface {
	Type() SplitType
	RatioIndexLeft() int
	RatioIndexRight() int
}
*/

type TreeRegions interface {
	Offset() *Vector
	Scale() float64
	Region(id NodeID) Region
}

type Tree struct {
	uniqueId     NodeID
	ratios       TreeRatios
	xyRatioIndex int
	zyRatioIndex int
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

func (tree *Tree) Parent(node *Node) *Node {
	return node.parent
}

func (tree *Tree) Ratios() Ratios {
	return tree.ratios.Ratios()
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
