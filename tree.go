package hambidgetree

// Serialize this like so
//
// ratios: (tree.ratios.ratios) // should marshal themselves with some kind of type flag, if strings, then expressions, otherwise floats
// xyRatioIndex:
// zyRatioIndex:
// root: id
// node: [id: id, split: hvd, index, left: id, right: id]

type NodeID int64

type Tree interface {
	RatioSource() RatioSource
	Node(id NodeID) Node
	Parent(id NodeID) Node
	Root() Node
	RatioIndexXY() int
	RatioIndexZY() int
}

type Branch interface {
	SplitType() SplitType
	Left() Node
	Right() Node
	LeftIndex() int
	RightIndex() int
}

type Node interface {
	ID() NodeID
	Branch() Branch
}

// This isn't part of ImmutableTree, it is returned by some builders
type Leaf interface {
	ID() NodeID
	RatioIndexXY() int
	RatioIndexZY() int
}

// An object used for iterating nodes within a tree or part of a tree
type NodeIterator struct {
	nodes []Node
}

// Creates a node iterator at the given node.
func NewNodeIterator(root Node) *NodeIterator {
	return &NodeIterator{
		nodes: []Node{root},
	}
}

// Whether the node iterator has more nodes to iterate.
func (it *NodeIterator) HasNext() bool {
	return len(it.nodes) > 0
}

// Returns the next node or nil
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
