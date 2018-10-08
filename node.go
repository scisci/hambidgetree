package hambidgetree

type NodeID int64

type ImmutableTree interface {
	Ratios() Ratios
	Node(id NodeID) ImmutableNode
	Parent(id NodeID) ImmutableNode
	Root() ImmutableNode
	RatioIndexXY() int
	RatioIndexZY() int
}

type ImmutableBranch interface {
	SplitType() SplitType
	Left() ImmutableNode
	Right() ImmutableNode
	LeftIndex() int
	RightIndex() int
}

type ImmutableNode interface {
	ID() NodeID
	Branch() ImmutableBranch
}

// This isn't part of ImmutableTree, it is returned by some builders
type ImmutableLeaf interface {
	ID() NodeID
	RatioIndexXY() int
	RatioIndexZY() int
}
