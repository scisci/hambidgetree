package neighbor

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/generators/grid"
	"sort"
	"testing"
)

type nodeListByID []htree.Node

func (a nodeListByID) Len() int           { return len(a) }
func (a nodeListByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodeListByID) Less(i, j int) bool { return a[i].ID() < a[j].ID() }

func arrayContentsEqual(arr1, arr2 []htree.Node) bool {
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
	tree := grid.New2D(2) // Create a tree with 4 squares
	leaves := htree.FindLeaves(tree)
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)

	epsilon := 0.0000001

	neighbors := getNeighbors(leaves, regionMap, epsilon)
	// Each one should have a neighbor
	if !arrayContentsEqual(leaves, neighbors) {
		t.Errorf("All items should be neighbors got %v, expected %v", neighbors, leaves)
	}

	// Remove one of the leaves
	leaves = leaves[:len(leaves)-1]
	neighbors = getNeighbors(leaves, regionMap, epsilon)
	if len(neighbors) != 2 {
		t.Errorf("Should only have two neighbors but got (%d)", len(neighbors))
	}

}
