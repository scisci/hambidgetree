package hambidgetree

import "testing"

func generateGridTree(levels int) *HambidgeTree {
	ratios := NewRatios([]float64{0.5, 1.0, 2.0})
	treeRatios := NewTreeRatios(ratios, 0.0000001)

	tree := NewHambidgeTree(treeRatios, 1)

	for i := 0; i < levels; i++ {
		leaves := tree.Leaves()
		for _, leaf := range leaves {
			var split Split
			if i&1 == 0 {
				split = NewVerticalSplit(0, 0)
			} else {
				split = NewHorizontalSplit(0, 0)
			}

			leaf.Divide(split)
		}
	}

	return tree
}

var strokeTests = []struct {
	Tree  *HambidgeTree
	Calls []GraphicsContextCall
}{
	{
		Tree: generateGridTree(1),
		Calls: []GraphicsContextCall{
			&GraphicsContextLine{0.5, 0.0, 0.5, 1.0},
			&GraphicsContextRect{0.0, 0.0, 1.0, 1.1},
		},
	},
}

func TestTreeStrokeRenderer(t *testing.T) {
	for i, test := range strokeTests {
		// Generates a 4x4 grid
		renderer := NewTreeStrokeRenderer()
		gc := NewGraphicsContextRecorder()

		renderer.Render(test.Tree, gc)

		if len(gc.Calls) != len(test.Calls) {
			t.Errorf("Tree stroke test %d failed, lengths don't match, expected %d, got %d", i, len(test.Calls), len(gc.Calls))
			continue
		}

		for c := 0; c < len(gc.Calls); c++ {
			if !test.Calls[c].Equals(gc.Calls[c]) {
				t.Errorf("Tree stroke test %d failed, call %d doesn't match, expected %v, got %v", i, c, test.Calls[c], gc.Calls[c])
			}
		}

	}

}
