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

type Node struct {
	tree   *Tree
	split  Split
	parent *Node
	left   *Node
	right  *Node
}

func NewNode(tree *Tree, parent *Node) *Node {
	node := &Node{
		tree:   tree,
		parent: parent,
	}

	return node
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

func (node *Node) Split() Split {
	return node.split
}

func (node *Node) IsLeft() bool {
	return node.parent != nil && node.parent.left == node
}

func (node *Node) RatioIndex() int {
	if node.parent != nil {
		if node.parent.left == node {
			return node.parent.split.LeftIndex()
		}

		return node.parent.split.RightIndex()
	}

	// No parent, must be the root
	return node.tree.ratioIndex
}

func (node *Node) Ratio() float64 {
	return node.tree.ratios.Ratios().At(node.RatioIndex())
}
