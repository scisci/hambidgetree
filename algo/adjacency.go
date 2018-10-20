package algo

import (
	htree "github.com/scisci/hambidgetree"
)

func BuildAdjacencyMatrix(tree htree.Tree, regionMap htree.RegionMap) map[htree.NodeID][]htree.Node {
	leaves := FindLeaves(tree)
	matrix := make(map[htree.NodeID][]htree.Node)

	for _, leaf := range leaves {
		neighbors := FindNeighbors(tree, leaf, regionMap)
		matrix[leaf.ID()] = neighbors
	}

	return matrix
}
