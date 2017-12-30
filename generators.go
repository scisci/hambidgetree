package hambidgetree

import "math/rand"
import "strconv"
import "errors"

type ParameterFormatType int

const (
	ParameterFormatTypeVerbose ParameterFormatType = 0
	ParameterFormatTypeConcise ParameterFormatType = 1
)

type TreeGenerator interface {
	Name() string
	Description() string
	Parameters(f ParameterFormatType) map[string]interface{}
	Generate() (*Tree, error)
}

type RandomBasicTreeGenerator struct {
	Ratios         TreeRatios
	ContainerRatio float64
	NumLeaves      int
	Seed           int64
	Is3D           bool
}

type leafSplits struct {
	leaf   *DimensionalNode
	splits []Split
}

func NewRandomBasicTreeGenerator(ratios TreeRatios, containerRatio float64, numLeaves int, seed int64) *RandomBasicTreeGenerator {
	return &RandomBasicTreeGenerator{
		Ratios:         ratios,
		ContainerRatio: containerRatio,
		NumLeaves:      numLeaves,
		Seed:           seed,
		Is3D:           false,
	}
}

func (gen *RandomBasicTreeGenerator) Name() string {
	return "Random Basic"
}

func (gen *RandomBasicTreeGenerator) Description() string {
	return "This algorithm can be set to generate a given number of leaves. It begins with a single leaf of a given ratio. Until it reaches the desired number of leaves, it selects a leaf at random and splits it."
}

func (gen *RandomBasicTreeGenerator) Parameters(f ParameterFormatType) map[string]interface{} {
	if f == ParameterFormatTypeConcise {
		return map[string]interface{}{
			"# Leaves": gen.NumLeaves,
			"Seed":     gen.Seed,
		}
	}

	return map[string]interface{}{
		"Ratios":           RatiosParameterString(gen.Ratios.Ratios()),
		"Container Ratio":  strconv.FormatFloat(gen.ContainerRatio, 'f', 4, 64),
		"Number of Leaves": gen.NumLeaves,
		"Random Seed":      gen.Seed,
		"3D":               gen.Is3D,
	}
}

func (gen *RandomBasicTreeGenerator) filterLeaves2D(leaf *DimensionalNode, complements Complements) *leafSplits {
	ratioIndex := leaf.tree.RatioIndex(leaf.Node, RatioPlaneXY)

	if len(complements[ratioIndex]) == 0 {
		return nil
	}

	return &leafSplits{
		leaf:   leaf,
		splits: complements[ratioIndex],
	}
}

func (gen *RandomBasicTreeGenerator) filterLeaves3D(leaf *DimensionalNode, complements Complements) *leafSplits {
	//ratioIndex := leaf.tree.RatioIndex(leaf.Node, RatioPlaneXY)

	if len(complements[leaf.RatioIndex()]) == 0 {
		return nil
	}

	return &leafSplits{
		leaf:   leaf,
		splits: complements[leaf.RatioIndex()],
	}
}

func (gen *RandomBasicTreeGenerator) Generate() (*Tree, error) {
	rand.Seed(gen.Seed)

	epsilon := CalculateRatiosEpsilon(gen.Ratios.Ratios())
	containerRatioIndex := FindClosestIndex(gen.Ratios.Ratios(), gen.ContainerRatio, epsilon)
	if containerRatioIndex < 0 {
		return nil, errors.New("Container ratio not found in list of ratios.")
	}

	complements := gen.Ratios.Complements()

	// Generate the container
	var tree *Tree
	if !gen.Is3D {
		tree = NewTree2D(gen.Ratios, containerRatioIndex)
	} else {
		tree = NewTree(gen.Ratios, containerRatioIndex, containerRatioIndex)
	}

	leafCount := 1

	leafDims := []*DimensionalNode{
		NewDimensionalNodeFromTree(tree, &Vector{0, 0, 0}, 1.0),
	}

	lookup := make(map[NodeID]*DimensionalNode)

	for {
		if leafCount >= gen.NumLeaves {
			break
		}

		it := NewDimensionalIteratorFromLeaves(leafDims)

		leafDims = leafDims[:0]
		for it.HasNext() {
			dimNode := it.Next()
			lookup[dimNode.Node.id] = dimNode
			if dimNode.IsLeaf() {
				leafDims = append(leafDims, dimNode)
			}
		}

		// Collect all splittable leaves
		var filteredLeaves []*leafSplits

		for _, leaf := range leafDims {

			if !gen.Is3D {
				if filteredLeaf := gen.filterLeaves2D(leaf, complements); filteredLeaf != nil {
					filteredLeaves = append(filteredLeaves, filteredLeaf)
				}
			} else {
				if filteredLeaf := gen.filterLeaves3D(leaf, complements); filteredLeaf != nil {
					filteredLeaves = append(filteredLeaves, filteredLeaf)
				}
			}
		}
		/*
			splittableLeaves := tree.FilterNodes(func(node *Node) bool {
				return node.IsLeaf() && len(complements[node.RatioIndex()]) > 0
			})
		*/

		if len(filteredLeaves) == 0 {
			return nil, errors.New("Unable to reach desired number of leaves (" + strconv.Itoa(gen.NumLeaves) + "), got " + strconv.Itoa(leafCount) + ".")
		}

		// Choose a leaf at random
		leafIndex := rand.Intn(len(filteredLeaves))
		filteredLeaf := filteredLeaves[leafIndex]
		splits := filteredLeaf.splits

		// Choose a random split
		splitIndex := rand.Intn(len(splits))
		split := splits[splitIndex]

		// Randomly invert the split (by default complements always have the smaller)
		// ratio on the left, but we want it to be evenly distributed.
		if rand.Int()&1 == 0 {
			split = split.Inverse()
		}

		filteredLeaf.leaf.Divide(split)
		leafCount += 1
	}

	return tree, nil
}
