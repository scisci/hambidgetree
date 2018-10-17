package algo

import (
	htree "github.com/scisci/hambidgetree"
)

// Go up the tree and select all 'other' leaves, then recursively visit any
// branches that intersect our leaf until we find leaves that intersect
func FindNeighbors(tree htree.Tree, node htree.Node, regionMap htree.RegionMap) []htree.Node {
	// TODO: performance test NodeDimension, if its too slow, just use a
	// DimensionalNode which has the hierarchy and the dimensions built in.
	dim := regionMap[node.ID()].Dimension()

	var neighbors []htree.Node

	epsilon := 0.0000001

	ref := node
	parent := tree.Parent(ref.ID())
	var stack []htree.Node
	for parent != nil {
		parentBranch := parent.Branch()

		if parentBranch == nil {
			panic("How can parent branch be nil!")
		}

		other := parentBranch.Left()
		if other == ref {
			other = parentBranch.Right()
		}
		stack = append(stack, other)
		ref = parent
		parent = tree.Parent(ref.ID())
	}

	for len(stack) > 0 {
		other := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		otherDim := regionMap[other.ID()].Dimension()

		dist := dim.DistanceSquared(otherDim)
		if dist > epsilon {
			// Not a neighbor
			continue
		}

		otherBranch := other.Branch()
		if otherBranch == nil {
			neighbors = append(neighbors, other)
		} else {
			stack = append(stack, otherBranch.Left(), otherBranch.Right())
		}
	}

	return neighbors
}
