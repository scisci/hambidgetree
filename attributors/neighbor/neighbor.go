package neighbor

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/algo"
	"github.com/scisci/hambidgetree/attributors"
	"math/rand"
)

var HasNeighborAttr = "hasNeighbor"
var HasNeighborValue = "true"

type HasNeighborAttributor struct {
	MaxMarks  int
	Dimension int
	Seed      int64
}

func NewHasNeighborAttributor(maxMarks int, dimension int, seed int64) *HasNeighborAttributor {
	return &HasNeighborAttributor{
		MaxMarks:  maxMarks,
		Dimension: dimension,
		Seed:      seed,
	}
}

func (attributor *HasNeighborAttributor) AddAttributes(tree htree.Tree, attrs *attributors.NodeAttributer) error {
	rand.Seed(attributor.Seed)
	epsilon := 0.0000001

	// Get a list of all the nodes
	leaves := algo.FindLeaves(tree)

	// Get the dimension list
	regionMap := htree.NewTreeRegionMap(tree, htree.Origin, htree.UnityScale)

	// Naive approach, just compare each leaf to each other leaf, could do better
	// with some sorting
	for count := 0; count < attributor.MaxMarks; count++ {
		// Find all the remaining leaves that still have neighbors
		leaves = getNeighbors(leaves, regionMap, epsilon)
		if len(leaves) < 2 {
			break
		}

		randomIndex := rand.Intn(len(leaves))
		attrs.SetAttribute(leaves[randomIndex].ID(), HasNeighborAttr, HasNeighborValue)
		leaves = append(leaves[:randomIndex], leaves[randomIndex+1:]...)
	}

	return nil
}
