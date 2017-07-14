package hambidgetree

import "testing"

var splitTests = []struct {
	split      Split
	typ        SplitType
	leftIndex  int
	rightIndex int
	valid      bool
}{
	{
		NewSplit(SplitTypeHorizontal, 0, 1),
		SplitTypeHorizontal,
		0,
		1,
		true,
	}, {
		NewSplit(SplitTypeVertical, 0, 1),
		SplitTypeVertical,
		0,
		1,
		true,
	}, {
		NewSplit(SplitTypeHorizontal, 0, 0),
		SplitTypeHorizontal,
		0,
		0,
		true,
	}, {
		NewSplit(SplitTypeHorizontal, 10, 30),
		SplitTypeHorizontal,
		10,
		30,
		true,
	}, {
		NewSplit(SplitTypeNone, 0, 0),
		SplitTypeNone,
		0,
		0,
		false,
	}, {
		NewSplit(SplitTypeHorizontal, -1, -1),
		SplitTypeHorizontal,
		-1,
		-1,
		false,
	}}

func TestNewSplit(t *testing.T) {
	for i, args := range splitTests {
		split := args.split
		if (args.typ == SplitTypeHorizontal && (!split.IsHorizontal() || split.IsVertical())) ||
			(args.typ == SplitTypeVertical && (!split.IsVertical() || split.IsHorizontal())) ||
			(args.typ == SplitTypeNone && (split.IsVertical() || split.IsHorizontal())) {
			t.Errorf("Split test %d failed, dimensionality is wrong")
		}

		if index := split.LeftIndex(); index != args.leftIndex {
			t.Errorf("Split test %d failed, left index should be %d, got %d", i, args.leftIndex, index)
		}

		if index := split.RightIndex(); index != args.rightIndex {
			t.Errorf("Split test %d failed, right index should be %d, got %d", i, args.rightIndex, index)
		}

		if valid := split.IsValid(); valid != args.valid {
			t.Errorf("Split test %d failed, valid should be %t, got %t", i, args.valid, valid)
		}
	}

}
