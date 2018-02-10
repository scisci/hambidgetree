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
	dimMap := htree.NewNodeDimensionMap(tree, htree.Origin, htree.UnityScale)
	leaves := tree.Leaves()
	leaf := leaves[0]

	neighbors, err := algo.FindNeighbors(leaf, dimMap)
	if err != nil {
		t.Errorf("Error finding neighbors %v", err)
	}

	if len(neighbors) != 3 {
		t.Errorf("Tree of 4 leaves should always have 3 neighbors, got %d", len(neighbors))
	}
}

func TestNeighbors3D(t *testing.T) {
	tree := grid.New3D(3)
	dimMap := htree.NewNodeDimensionMap(tree, htree.Origin, htree.UnityScale)
	leaves := tree.Leaves()
	leaf := leaves[0]

	neighbors, err := algo.FindNeighbors(leaf, dimMap)
	if err != nil {
		t.Errorf("Error finding neighbors %v", err)
	}

	if len(neighbors) != 7 {
		t.Errorf("3d Tree of %d leaves should always have 7 neighbors, got %d", len(leaves), len(neighbors))
	}
}

func TestNeighbors3DMeasured(t *testing.T) {
	ratios := golden.Ratios()
	treeRatios := htree.NewTreeRatios(ratios, 0.0000001)
	numLeaves := 10
	gen := randombasic.New3D(treeRatios, 1, 1, numLeaves, 543543)
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	dimMap := htree.NewNodeDimensionMap(tree, htree.Origin, htree.UnityScale)
	leaves := tree.Leaves()

	for _, leaf := range leaves {
		neighbors, err := algo.FindNeighbors(leaf, dimMap)
		if err != nil {
			t.Errorf("Error finding neighbors %v", err)
		}

		dim, err := dimMap.Dimension(leaf.ID())
		if err != nil {
			t.Errorf("Error finding dimension of leaf %v", err)
		}

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

			otherDim, err := dimMap.Dimension(other.ID())
			if err != nil {
				t.Errorf("Error getting dimension of neighbor %v", err)
			}

			isCalcNeighbor := dim.DistanceSquared(otherDim) < 0.0000001
			if isCalcNeighbor != isNeighbor {
				t.Errorf("Neighbor incorrect, expected %t, got %t", isNeighbor, isCalcNeighbor)
			}
		}
	}
}

func TestAdjacencyMatrix(t *testing.T) {
	tree := grid.New2D(2)
	dimensionLookup := htree.NewNodeDimensionMap(tree, htree.Origin, htree.UnityScale)
	matrix, err := algo.BuildAdjacencyMatrix(tree, dimensionLookup)
	if err != nil {
		t.Errorf("Error building adjacency matrix %v", err)
	}

	leaves := tree.Leaves()
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
