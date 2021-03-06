package algo_test

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/algo"
	"github.com/scisci/hambidgetree/generators/grid"
	"github.com/scisci/hambidgetree/generators/randombasic"
	"github.com/scisci/hambidgetree/golden"
	"testing"
)

func TestNeighbors2D(t *testing.T) {
	tree := grid.New2D(2)
	leaves := algo.FindLeaves(tree)
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)

	leaf := leaves[0]

	neighbors := algo.FindNeighbors(tree, leaf, regionMap)
	if len(neighbors) != 3 {
		t.Errorf("Tree of 4 leaves should always have 3 neighbors, got %d", len(neighbors))
	}
}

func TestNeighbors3D(t *testing.T) {
	tree := grid.New3D(3)
	leaves := algo.FindLeaves(tree)
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)

	leaf := leaves[0]

	neighbors := algo.FindNeighbors(tree, leaf, regionMap)

	if len(neighbors) != 7 {
		t.Errorf("3d Tree of %d leaves should always have 7 neighbors, got %d", len(leaves), len(neighbors))
	}
}

func TestNeighbors3DMeasured(t *testing.T) {
	ratioSource := golden.RatioSource()
	numLeaves := 10
	gen, err := randombasic.New3D(ratioSource, 1, 1, numLeaves, 543543)
	if err != nil {
		t.Errorf("Error creating tree %v", err)
	}
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	leaves := algo.FindLeaves(tree)
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)

	for _, leaf := range leaves {
		neighbors := algo.FindNeighbors(tree, leaf, regionMap)

		dim := regionMap[leaf.ID()].AlignedBox()

		// Make sure that each leaf that is not in neighbors is actually not a
		// neighbor.
		for _, other := range leaves {
			if other == leaf {
				continue
			}

			isNeighbor := false
			for _, n := range neighbors {
				if n == other {
					isNeighbor = true
					break
				}
			}

			otherDim := regionMap[other.ID()].AlignedBox()
			isCalcNeighbor := dim.DistanceSquared(otherDim) < 0.0000001
			if isCalcNeighbor != isNeighbor {
				t.Errorf("Neighbor incorrect, expected %t, got %t", isNeighbor, isCalcNeighbor)
			}
		}
	}
}

func TestAdjacencyMatrix(t *testing.T) {
	tree := grid.New2D(2)
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)
	matrix := algo.BuildAdjacencyMatrix(tree, regionMap)
	leaves := algo.FindLeaves(tree)

	if len(matrix) != len(leaves) {
		t.Errorf("Should have %d items in matrix, got %d", len(leaves), len(matrix))
	}

	for nodeID, list := range matrix {
		if len(list) != 3 {
			t.Errorf("Each node should have 3 neighbors, got %d", len(list))
		}

		for _, otherNode := range list {
			if otherNode.ID() == nodeID {
				t.Errorf("Node has self as neighbor!")
			}
		}
	}
}
