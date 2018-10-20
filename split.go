package hambidgetree

type SplitType int

const SplitTypeHorizontal SplitType = 1 // Split along Y axis
const SplitTypeVertical SplitType = 2   // Split along X axis
const SplitTypeDepth SplitType = 3      // Split along Z axis

type Split struct {
	typ        SplitType
	leftIndex  int
	rightIndex int
}

func (s Split) Type() SplitType {
	return s.typ
}

func (s Split) LeftIndex() int {
	return s.leftIndex
}

func (s Split) RightIndex() int {
	return s.rightIndex
}

// Index arguments are 0 based, but internally we start at 1. This allows the
// default value of the interface to be obviously invalid.
func NewSplit(typ SplitType, leftIndex, rightIndex int) Split {
	return Split{
		typ:        typ,
		leftIndex:  leftIndex,
		rightIndex: rightIndex,
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

func NewInvertedSplit(s Split) Split {
	return Split{
		typ:        s.typ,
		leftIndex:  s.rightIndex,
		rightIndex: s.leftIndex,
	}
}
