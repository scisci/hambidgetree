package hambidgetree

func BuildAdjacencyMatrix(tree *Tree, dimensionLookup NodeDimensions) (map[NodeID][]*Node, error) {

	leaves := tree.Leaves()

	matrix := make(map[NodeID][]*Node)

	for _, leaf := range leaves {
		neighbors, err := FindNeighbors(leaf, dimensionLookup)
		if err != nil {
			return nil, err
		}
		matrix[leaf.id] = neighbors
	}
	return matrix, nil
}

// Go up the tree and select all 'other' leaves, then recursively visit any
// branches that intersect our leaf until we find leaves that intersect
func FindNeighbors(leaf *Node, dimensionLookup NodeDimensions) ([]*Node, error) {
	// TODO: performance test NodeDimension, if its too slow, just use a
	// DimensionalNode which has the hierarchy and the dimensions built in.
	dim, err := dimensionLookup.Dimension(leaf.ID())
	if err != nil {
		return nil, err
	}

	var neighbors []*Node

	epsilon := 0.0000001

	ref := leaf
	parent := ref.parent
	var stack []*Node
	for parent != nil {
		other := parent.left
		if other == ref {
			other = parent.right
		}
		stack = append(stack, other)
		ref = parent
		parent = ref.parent
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
