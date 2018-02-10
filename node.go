package hambidgetree

type NodeID int64

type RatioPlane int

const RatioPlaneXY = 1
const RatioPlaneZY = 2

type Node struct {
	tree   *Tree
	id     NodeID
	split  Split
	parent *Node
	left   *Node
	right  *Node
}

func NewNode(tree *Tree, parent *Node) *Node {
	node := &Node{
		id:     tree.nextNodeId(),
		tree:   tree,
		parent: parent,
	}

	return node
}

func (node *Node) ID() NodeID {
	return node.id
}

func (node *Node) IsLeaf() bool {
	return node.split == nil
}

func (node *Node) Divide(split Split) {
	if !node.IsLeaf() {
		panic("Node can't be split, not a leaf")
		return
	}

	node.split = split
	node.left = NewNode(node.tree, node)
	node.right = NewNode(node.tree, node)
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) Left() *Node {
	return node.left
}

func (node *Node) Right() *Node {
	return node.right
}

func (node *Node) Split() Split {
	return node.split
}

func (node *Node) IsLeft() bool {
	return node.parent != nil && node.parent.left == node
}

// Search up the tree looking for the last split that would affect the required
// axis.
func (node *Node) RatioIndex() int {
	if node.parent == nil {
		panic("Can't get ratio index this is the root!")
	}

	if node.parent.left == node {
		return node.parent.split.LeftIndex()
	}

	return node.parent.split.RightIndex()
}

func (node *Node) Ratio() float64 {
	return node.tree.ratios.Ratios().At(node.RatioIndex())
}
