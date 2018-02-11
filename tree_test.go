package hambidgetree_test

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/generators/grid"
	"testing"
)

func TestTreeGrid2D(t *testing.T) {
	tree := grid.New2D(4)
	leaves := htree.FindLeaves(tree)
	regionMap := htree.NewNodeRegionMap(tree, htree.Origin, htree.UnityScale)

	baseDim := regionMap.Region(tree.Root().ID()).Dimension()
	if baseDim.Width() != 1.0 {
		t.Errorf("Tree should have width 1, got %f", baseDim.Width())
	}

	if baseDim.Height() != 1.0 {
		t.Errorf("Tree should have height 1, got %f", baseDim.Height())
	}

	for _, leaf := range leaves {
		dim := regionMap.Region(leaf.ID()).Dimension()
		if dim.Width() != 0.25 {
			t.Errorf("Tree cell width should be 0.25, got %f", dim.Width())
		}
		if dim.Height() != 0.25 {
			t.Errorf("Tree cell height should be 0.25, got %f", dim.Height())
		}
		if dim.Depth() != 0.0 {
			t.Errorf("Tree cell depth should be 0.0, got %f", dim.Depth())
		}
	}
}

var grid3DTests = []struct {
	levels int
	width  float64
	height float64
	depth  float64
}{
	{
		levels: 6,
		width:  0.25,
		height: 0.25,
		depth:  0.25,
	},
	{
		levels: 1,
		width:  0.5,
		height: 1.0,
		depth:  1.0,
	},
	{
		levels: 2,
		width:  0.5,
		height: 0.5,
		depth:  1.0,
	},
	{
		levels: 3,
		width:  0.5,
		height: 0.5,
		depth:  0.5,
	},
}

func TestTreeGrid3D(t *testing.T) {
	for _, test := range grid3DTests {
		tree := grid.New3D(test.levels)
		leaves := htree.FindLeaves(tree)
		regionMap := htree.NewNodeRegionMap(tree, htree.Origin, htree.UnityScale)

		baseDim := regionMap.Region(tree.Root().ID()).Dimension()
		if baseDim.Width() != 1.0 {
			t.Errorf("Tree should have width 1, got %f", baseDim.Width())
		}

		if baseDim.Height() != 1.0 {
			t.Errorf("Tree should have height 1, got %f", baseDim.Height())
		}

		if baseDim.Depth() != 1.0 {
			t.Errorf("Tree should have depth 1, got %f", baseDim.Depth())
		}

		for _, leaf := range leaves {
			dim := regionMap.Region(leaf.ID()).Dimension()
			if dim.Width() != test.width {
				t.Errorf("Tree cell width should be %f got %f", test.width, dim.Width())
			}
			if dim.Height() != test.height {
				t.Errorf("Tree cell height should be %f, got %f", test.height, dim.Height())
			}
			if dim.Depth() != test.depth {
				t.Errorf("Tree cell depth should be %f, got %f", test.depth, dim.Depth())
			}
		}
	}
}
