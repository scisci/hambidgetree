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

type TreeRegions interface {
	Offset() *Vector
	Scale() float64
	Region(id NodeID) Region
}
