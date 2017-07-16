package hambidgetree

import "math/rand"
import "strconv"
import "errors"

type TreeGenerator interface {
	Name() string
	Description() string
	Parameters() map[string]interface{}
	Generate() (*Tree, error)
}

type RandomBasicTreeGenerator struct {
	Ratios         TreeRatios
	ContainerRatio float64
	NumLeaves      int
	Seed           int64
}

func NewRandomBasicTreeGenerator(ratios TreeRatios, containerRatio float64, numLeaves int, seed int64) *RandomBasicTreeGenerator {
	return &RandomBasicTreeGenerator{
		Ratios:         ratios,
		ContainerRatio: containerRatio,
		NumLeaves:      numLeaves,
		Seed:           seed,
	}
}

func (gen *RandomBasicTreeGenerator) Name() string {
	return "Random Basic"
}

func (gen *RandomBasicTreeGenerator) Description() string {
	return "This algorithm begins with a single leaf of a given ratio. Until it reaches the desired number of leaves, it selects a leaf at random and splits it."
}

func (gen *RandomBasicTreeGenerator) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"Ratios":           RatiosParameterString(gen.Ratios.Ratios()),
		"ContainerRatio":   strconv.FormatFloat(gen.ContainerRatio, 'f', 4, 64),
		"Number of Leaves": gen.NumLeaves,
		"Random Seed":      gen.Seed,
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
	tree := NewTree(gen.Ratios, containerRatioIndex)

	leafCount := 1

	for {
		if leafCount >= gen.NumLeaves {
			break
		}

		// Collect all splittable leaves
		splittableLeaves := tree.FilterNodes(func(node *Node) bool {
			return node.IsLeaf() && len(complements[node.RatioIndex()]) > 0
		})

		if len(splittableLeaves) == 0 {
			return nil, errors.New("Unable to reach desired number of leaves (" + strconv.Itoa(gen.NumLeaves) + "), got " + strconv.Itoa(leafCount) + ".")
		}

		// Choose a leaf at random
		leafIndex := rand.Intn(len(splittableLeaves))
		leaf := splittableLeaves[leafIndex]
		splits := complements[leaf.RatioIndex()]

		// Choose a random split
		splitIndex := rand.Intn(len(splits))
		split := splits[splitIndex]

		// Randomly invert the split (by default complements always have the smaller)
		// ratio on the left, but we want it to be evenly distributed.
		if rand.Int()&1 == 0 {
			split = split.Inverse()
		}

		splittableLeaves[leafIndex].Divide(split)
		leafCount += 1
	}

	return tree, nil
}
