package hambidgetree

import "math/rand"

type TreeAttributor interface {
	Name() string
	Description() string
	Parameters(f ParameterFormatType) map[string]interface{}
	AddAttributes(tree *Tree, attrs *NodeAttributer) error
}

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

func (attributor *HasNeighborAttributor) Name() string {
	return "HasNeighbor"
}

func (attributor *HasNeighborAttributor) Description() string {
	return "This attributor finds all of the leaves that have at least one neighbor touching on either side, then chooses a random one and marks it. Once marked, the leaf is considered 'deleted.'"
}

func (attributor *HasNeighborAttributor) Parameters(f ParameterFormatType) map[string]interface{} {
	return map[string]interface{}{
		"Max Marks": attributor.MaxMarks,
		"Dimension": attributor.Dimension,
		"Seed":      attributor.Seed,
	}
}

func (attributor *HasNeighborAttributor) AddAttributes(tree *Tree, attrs *NodeAttributer) error {
	rand.Seed(attributor.Seed)
	epsilon := 0.0000001

	// Get a list of all the nodes
	leaves := tree.Leaves()

	// Get the dimension list
	nodeDimMap := NewNodeDimensionMap(tree.root, 0, 0, 1.0)

	// Naive approach, just compare each leaf to each other leaf, could do better
	// with some sorting
	for count := 0; count < attributor.MaxMarks; count++ {
		// Find all the remaining leaves that still have neighbors
		leaves, err := getNeighbors(leaves, nodeDimMap, epsilon)
		if err != nil {
			return err
		}

		if len(leaves) < 2 {
			break
		}

		randomIndex := rand.Intn(len(leaves))
		attrs.SetAttribute(leaves[randomIndex], HasNeighborAttr, HasNeighborValue)
		leaves = append(leaves[:randomIndex], leaves[randomIndex+1:]...)
	}

	return nil
}

func getNeighbors(leaves []*Node, nodeDimMap NodeDimensions, epsilon float64) ([]*Node, error) {
	for i := 0; i < len(leaves); {
		hasNeighbor := false
		dim, err := nodeDimMap.Dimension(leaves[i])
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(leaves); j++ {
			if j == i {
				continue
			}

			dim2, err := nodeDimMap.Dimension(leaves[j])
			if err != nil {
				return nil, err
			}

			if !dim.IntersectLeft(dim2, epsilon).Empty() ||
				!dim.IntersectRight(dim2, epsilon).Empty() {
				hasNeighbor = true
				break
			}
		}

		if !hasNeighbor {
			leaves = append(leaves[:i], leaves[i+1:]...) // Remove the candidate, keep i
		} else {
			i++ // Has neighbor so move on
		}
	}

	return leaves, nil
}
