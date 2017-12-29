package hambidgetree

import "testing"
import "sort"

type nodeListByID []*Node

func (a nodeListByID) Len() int           { return len(a) }
func (a nodeListByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodeListByID) Less(i, j int) bool { return a[i].id < a[j].id }

func arrayContentsEqual(arr1, arr2 []*Node) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	sort.Sort(nodeListByID(arr1))
	sort.Sort(nodeListByID(arr2))

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func TestGetNeighbors(t *testing.T) {
	tree := NewGridTree(2) // Create a tree with 4 squares
	leaves := tree.Leaves()
	epsilon := 0.0000001
	nodeDimMap := NewNodeDimensionMap(tree.root, 0, 0, 1.0)
	neighbors, err := getNeighbors(leaves, nodeDimMap, epsilon)
	if err != nil {
		t.Errorf("Error getting neighbors (%v)", err)
	}

	// Each one should have a neighbor
	if !arrayContentsEqual(leaves, neighbors) {
		t.Errorf("All items should be neighbors got %v, expected %v", neighbors, leaves)
	}

	// Remove one of the leaves
	leaves = leaves[:len(leaves)-1]
	neighbors, err = getNeighbors(leaves, nodeDimMap, epsilon)
	if err != nil {
		t.Errorf("Error getting neighbors (%v)", err)
	}

	if len(neighbors) != 2 {
		t.Errorf("Should only have two neighbors but got (%d)", len(neighbors))
	}

}