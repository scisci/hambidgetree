package hambidgetree

type HambidgeTree struct {
	ratios     TreeRatios
	ratioIndex int
	scale      float64
	root       *HambidgeTreeNode
}

/*
type NodesByLex []*HambidgeTreeNode

func (nodes NodesByLex) Len() int      { return len(nodes) }
func (nodes NodesByLex) Swap(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] }
func (nodes NodesByLex) Less(i, j int) bool {
	if nodes[i].Y() < nodes[j].Y() {
		return true
	}

	if nodes[j].Y() < nodes[i].Y() {
		return false
	}

	if nodes[i].X() < nodes[j].X() {
		return true
	}

	return false
}
*/
func NewHambidgeTree(ratios TreeRatios, ratioIndex int) *HambidgeTree {
	if ratioIndex < 0 || ratioIndex >= ratios.Ratios().Len() {
		panic("Invalid ratio index")
	}

	tree := &HambidgeTree{
		ratios:     ratios,
		ratioIndex: ratioIndex,
		scale:      1.0,
	}

	tree.root = NewHambidgeTreeNode(tree, nil)
	return tree
}

func (tree *HambidgeTree) Leaves() []*HambidgeTreeNode {
	return tree.FilterNodes(func(node *HambidgeTreeNode) bool {
		return node.IsLeaf()
	})
}

func (tree *HambidgeTree) FilterNodes(filter func(*HambidgeTreeNode) bool) []*HambidgeTreeNode {
	var nodes []*HambidgeTreeNode
	it := NewHambidgeTreeNodeIterator(tree.root)
	for it.HasNext() {
		node := it.Next()
		if filter(node) {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree *HambidgeTree) Ratio() float64 {
	return tree.ratios.Ratios().At(tree.ratioIndex)
}

type HambidgeTreeNodeIterator struct {
	nodes []*HambidgeTreeNode
}

func NewHambidgeTreeNodeIterator(root *HambidgeTreeNode) *HambidgeTreeNodeIterator {
	return &HambidgeTreeNodeIterator{
		nodes: []*HambidgeTreeNode{root},
	}
}

func (it *HambidgeTreeNodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

func (it *HambidgeTreeNodeIterator) Next() *HambidgeTreeNode {
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

/*
type HambidgeTreeNodeStats struct {
	width  float64
	height float64
	x      float64
	y      float64
}
*/
type HambidgeTreeNode struct {
	tree   *HambidgeTree
	split  Split
	parent *HambidgeTreeNode
	left   *HambidgeTreeNode
	right  *HambidgeTreeNode

	// Cached information
	//cachedStats HambidgeTreeNodeStats
}

func NewHambidgeTreeNode(tree *HambidgeTree, parent *HambidgeTreeNode) *HambidgeTreeNode {
	node := &HambidgeTreeNode{
		tree:   tree,
		parent: parent,
	}

	//node.cachedStats.x = -1.0
	//node.cachedStats.y = -1.0
	return node
}

func (node *HambidgeTreeNode) IsLeaf() bool {
	return node.split == nil
}

func (node *HambidgeTreeNode) Divide(split Split) {
	if !node.IsLeaf() {
		panic("Node can't be split, not a leaf")
		return
	}

	node.split = split
	node.left = NewHambidgeTreeNode(node.tree, node)
	node.right = NewHambidgeTreeNode(node.tree, node)
}

func (node *HambidgeTreeNode) Split() Split {
	return node.split
}

func (node *HambidgeTreeNode) IsLeft() bool {
	return node.parent != nil && node.parent.left == node
}

func (node *HambidgeTreeNode) RatioIndex() int {
	if node.parent != nil {
		if node.parent.left == node {
			return node.parent.split.LeftIndex()
		}

		return node.parent.split.RightIndex()
	}

	// No parent, must be the root
	return node.tree.ratioIndex
}

func (node *HambidgeTreeNode) Ratio() float64 {
	return node.tree.ratios.Ratios().At(node.RatioIndex())
}

/*

func (node *HambidgeTreeNode) Area() float64 {
	return node.Width() * node.Height()
}

func (node *HambidgeTreeNode) Width() float64 {
	if node.cachedStats.width <= 0.0 {
		if node.parent != nil {
			if node.parent.Split().IsHorizontal() {
				node.cachedStats.width = node.parent.Width()
			} else {
				node.cachedStats.width = node.parent.Width() * node.Ratio() / node.parent.Ratio()
			}
		} else {
			node.cachedStats.width = node.tree.Ratio() * node.tree.scale
		}
	}

	return node.cachedStats.width
}

func (node *HambidgeTreeNode) Height() float64 {
	if node.cachedStats.height <= 0.0 {
		if node.parent != nil {
			if node.parent.Split().IsVertical() {
				node.cachedStats.height = node.parent.Height()
			} else {
				node.cachedStats.height = node.parent.Height() * node.parent.Ratio() / node.Ratio()
			}
		} else {
			node.cachedStats.height = node.tree.scale
		}
	}

	return node.cachedStats.height
}

func (node *HambidgeTreeNode) X() float64 {
	if node.cachedStats.x < 0.0 {
		if node.parent != nil {
			if node.parent.Split().IsHorizontal() || node.IsLeft() {
				node.cachedStats.x = node.parent.X()
			} else {
				node.cachedStats.x = node.parent.X() + node.parent.left.Width()
			}
		} else {
			node.cachedStats.x = 0.0
		}
	}

	return node.cachedStats.x
}

func (node *HambidgeTreeNode) Y() float64 {
	if node.cachedStats.y < 0.0 {
		if node.parent != nil {
			if node.parent.Split().IsVertical() || node.IsLeft() {
				node.cachedStats.y = node.parent.Y()
			} else {
				node.cachedStats.y = node.parent.Y() + node.parent.left.Height()
			}
		} else {
			node.cachedStats.y = 0.0
		}
	}

	return node.cachedStats.y
}
*/
