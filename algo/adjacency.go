package algo

import (
	htree "github.com/scisci/hambidgetree"
)

func BuildAdjacencyMatrix(tree *htree.Tree, dimensionLookup htree.NodeDimensions) (map[htree.NodeID][]*htree.Node, error) {
	leaves := tree.Leaves()

	matrix := make(map[htree.NodeID][]*htree.Node)

	for _, leaf := range leaves {
		neighbors, err := FindNeighbors(leaf, dimensionLookup)
		if err != nil {
			return nil, err
		}
		matrix[leaf.ID()] = neighbors
	}
	return matrix, nil
}
