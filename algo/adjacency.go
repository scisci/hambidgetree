package algo

import (
	htree "github.com/scisci/hambidgetree"
)

func BuildAdjacencyMatrix(tree htree.ImmutableTree, regionLookup htree.TreeRegions) map[htree.NodeID][]htree.ImmutableNode {
	leaves := htree.FindLeaves(tree)
	matrix := make(map[htree.NodeID][]htree.ImmutableNode)

	for _, leaf := range leaves {
		neighbors := FindNeighbors(tree, leaf, regionLookup)
		matrix[leaf.ID()] = neighbors
	}

	return matrix
}
