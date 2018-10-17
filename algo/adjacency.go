package algo

import (
	htree "github.com/scisci/hambidgetree"
)

func BuildAdjacencyMatrix(tree htree.Tree, regionLookup htree.TreeRegions) map[htree.NodeID][]htree.Node {
	leaves := htree.FindLeaves(tree)
	matrix := make(map[htree.NodeID][]htree.Node)

	for _, leaf := range leaves {
		neighbors := FindNeighbors(tree, leaf, regionLookup)
		matrix[leaf.ID()] = neighbors
	}

	return matrix
}
