package hambidgetree

import "strconv"

type SplitType int

const SplitTypeHorizontal SplitType = 1 // Split along Y axis
const SplitTypeVertical SplitType = 2   // Split along X axis
const SplitTypeDepth SplitType = 3      // Split along Z axis

type Split interface {
	LeftIndex() int
	RightIndex() int
	IsHorizontal() bool
	IsVertical() bool
	IsDepth() bool
	String() string
	Type() SplitType
	Inverse() Split
}

type split struct {
	typ        SplitType
	leftIndex  int // Indices start at 1, so 0 is invalid
	rightIndex int // Indices start at 1, so 0 is invalid
}

func (s *split) Type() SplitType {
	return s.typ
}

func (s *split) Inverse() Split {
	return &split{
		typ:        s.typ,
		leftIndex:  s.rightIndex,
		rightIndex: s.leftIndex,
	}
}

func (s *split) LeftIndex() int {
	return s.leftIndex - 1
}

func (s *split) RightIndex() int {
	return s.rightIndex - 1
}

func (s *split) IsHorizontal() bool {
	return s.typ == SplitTypeHorizontal
}

func (s *split) IsVertical() bool {
	return s.typ == SplitTypeVertical
}

func (s *split) IsDepth() bool {
	return s.typ == SplitTypeDepth
}

func (s *split) String() string {
	str := "Split{"
	switch s.typ {
	case SplitTypeHorizontal:
		str = str + "h"
	case SplitTypeVertical:
		str = str + "v"
	case SplitTypeDepth:
		str = str + "d"
	default:
		str = str + "?"
	}

	str = str + "," + strconv.Itoa(s.leftIndex-1)
	str = str + "," + strconv.Itoa(s.rightIndex-1)
	str = str + "}"
	return str
}

// Index arguments are 0 based, but internally we start at 1. This allows the
// default value of the interface to be obviously invalid.
func NewSplit(typ SplitType, leftIndex, rightIndex int) Split {
	return &split{
		typ:        typ,
		leftIndex:  leftIndex + 1,
		rightIndex: rightIndex + 1,
	}
}

func NewVerticalSplit(leftIndex, rightIndex int) Split {
	return NewSplit(SplitTypeVertical, leftIndex, rightIndex)
}

func NewHorizontalSplit(leftIndex, rightIndex int) Split {
	return NewSplit(SplitTypeHorizontal, leftIndex, rightIndex)
}

func NewDepthSplit(leftIndex, rightIndex int) Split {
	return NewSplit(SplitTypeDepth, leftIndex, rightIndex)
}
