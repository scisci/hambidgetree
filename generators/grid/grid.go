package grid

import (
	htree "github.com/scisci/hambidgetree"
)

// Useful for testing, each level splits the tree in half, starting with a
// square
// Even levels split all leaves in half using a vertical line
// Odd levels split those splits in half with a horizontal line creating squares again
func New2D(levels int) *htree.Tree {
	ratios := htree.NewRatios([]float64{0.5, 1.0})
	treeRatios := htree.NewTreeRatios(ratios, 0.0000001)

	tree := htree.NewTree2D(treeRatios, 1)

	for i := 0; i < levels; i++ {
		leaves := tree.Leaves()
		for _, leaf := range leaves {
			var split htree.Split
			if i&1 == 0 {
				split = htree.NewVerticalSplit(0, 0)
			} else {
				split = htree.NewHorizontalSplit(1, 1)
			}

			leaf.Divide(split)
		}
	}

	return tree
}

func New3D(levels int) *htree.Tree {
	ratios := htree.NewRatios([]float64{0.5, 1.0})
	treeRatios := htree.NewTreeRatios(ratios, 0.0000001)

	tree := htree.NewTree(treeRatios, 1, 1)

	for i := 0; i < levels; i++ {
		leaves := tree.Leaves()
		for _, leaf := range leaves {
			var split htree.Split
			index := i % 3
			switch index {
			case 0:
				split = htree.NewVerticalSplit(0, 0)
			case 1:
				split = htree.NewHorizontalSplit(1, 1)
			case 2:
				split = htree.NewDepthSplit(1, 1)
			}

			leaf.Divide(split)
		}
	}

	return tree
}
