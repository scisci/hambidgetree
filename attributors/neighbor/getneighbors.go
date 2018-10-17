package neighbor

import (
	htree "github.com/scisci/hambidgetree"
)

func getNeighbors(leaves []htree.Node, regionMap htree.RegionMap, epsilon float64) []htree.Node {
	var candidates []htree.Node
	for i := 0; i < len(leaves); i++ {
		hasNeighbor := false
		dim := regionMap[leaves[i].ID()].Dimension()

		for j := 0; j < len(leaves); j++ {
			if j == i {
				continue
			}
			dim2 := regionMap[leaves[j].ID()].Dimension()

			leftExtent := dim.IntersectLeft(dim2, epsilon)
			rightExtent := dim.IntersectRight(dim2, epsilon)
			if !leftExtent.Empty() {
				//fmt.Printf("L/I %v, %v, %v\n", leftExtent, dim, dim2)
				hasNeighbor = true
			} else if !rightExtent.Empty() {
				//fmt.Printf("R/I %v, %v, %v\n", rightExtent, dim, dim2)
				hasNeighbor = true
				break
			}
		}

		if hasNeighbor {
			candidates = append(candidates, leaves[i]) // Remove the candidate, keep i
		}
	}

	return candidates
}
