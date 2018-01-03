package hambidgetree

import "testing"

func TestNeighbors2D(t *testing.T) {
	tree := NewGridTree2D(2)
	dimMap := NewNodeDimensionMap(tree, NewVector(0, 0, 0), 1.0)
	leaves := tree.Leaves()
	leaf := leaves[0]

	neighbors, err := FindNeighbors(leaf, dimMap)
	if err != nil {
		t.Errorf("Error finding neighbors %v", err)
	}

	if len(neighbors) != 3 {
		t.Errorf("Tree of 4 leaves should always have 3 neighbors, got %d", len(neighbors))
	}
}

func TestNeighbors3D(t *testing.T) {
	for l := 3; l < 10; l++ {
		tree := NewGridTree3D(l)
		dimMap := NewNodeDimensionMap(tree, NewVector(0, 0, 0), 1.0)
		leaves := tree.Leaves()
		leaf := leaves[0]

		neighbors, err := FindNeighbors(leaf, dimMap)
		if err != nil {
			t.Errorf("Error finding neighbors %v", err)
		}

		if len(neighbors) != 7 {
			t.Errorf("3d Tree of %d leaves should always have 7 neighbors, got %d", l, len(neighbors))
		}
	}
}
