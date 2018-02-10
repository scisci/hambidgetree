package algo

import (
	htree "github.com/scisci/hambidgetree"
)

// Go up the tree and select all 'other' leaves, then recursively visit any
// branches that intersect our leaf until we find leaves that intersect
func FindNeighbors(leaf *htree.Node, dimensionLookup htree.NodeDimensions) ([]*htree.Node, error) {
	// TODO: performance test NodeDimension, if its too slow, just use a
	// DimensionalNode which has the hierarchy and the dimensions built in.
	dim, err := dimensionLookup.Dimension(leaf.ID())
	if err != nil {
		return nil, err
	}

	var neighbors []*htree.Node

	epsilon := 0.0000001

	ref := leaf
	parent := ref.Parent()
	var stack []*htree.Node
	for parent != nil {
		other := parent.Left()
		if other == ref {
			other = parent.Right()
		}
		stack = append(stack, other)
		ref = parent
		parent = ref.Parent()
	}

	for len(stack) > 0 {
		branch := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		otherDim, err := dimensionLookup.Dimension(branch.ID())
		if err != nil {
			return nil, err
		}

		dist := dim.DistanceSquared(otherDim)
		if dist > epsilon {
			// Not a neighbor
			continue
		}

		if branch.IsLeaf() {
			neighbors = append(neighbors, branch)
		} else {
			stack = append(stack, branch.Left(), branch.Right())
		}
	}

	return neighbors, nil
}
