package hambidgetree

import "testing"

var splitTests = []struct {
	split      Split
	typ        SplitType
	leftIndex  int
	rightIndex int
}{
	{
		NewSplit(SplitTypeHorizontal, 0, 1),
		SplitTypeHorizontal,
		0,
		1,
	}, {
		NewSplit(SplitTypeVertical, 0, 1),
		SplitTypeVertical,
		0,
		1,
	}, {
		NewSplit(SplitTypeDepth, 0, 1),
		SplitTypeDepth,
		0,
		1,
	}, {
		NewSplit(SplitTypeHorizontal, 0, 0),
		SplitTypeHorizontal,
		0,
		0,
	}, {
		NewSplit(SplitTypeHorizontal, 10, 30),
		SplitTypeHorizontal,
		10,
		30,
	}, {
		NewSplit(SplitTypeHorizontal, -1, -1),
		SplitTypeHorizontal,
		-1,
		-1,
	}}

func TestNewSplit(t *testing.T) {
	for i, args := range splitTests {
		split := args.split

		if index := split.leftIndex; index != args.leftIndex {
			t.Errorf("Split test %d failed, left index should be %d, got %d", i, args.leftIndex, index)
		}

		if index := split.rightIndex; index != args.rightIndex {
			t.Errorf("Split test %d failed, right index should be %d, got %d", i, args.rightIndex, index)
		}
	}

}
