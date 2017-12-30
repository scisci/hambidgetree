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
	Ratios    TreeRatios
	XYRatio   float64
	ZYRatio   float64
	NumLeaves int
	Seed      int64
}

type leafSplits struct {
	leaf   *DimensionalNode
	splits []Split
}

func NewRandomBasicTreeGenerator(ratios TreeRatios, containerRatio float64, numLeaves int, seed int64) *RandomBasicTreeGenerator {
	return &RandomBasicTreeGenerator{
		Ratios:    ratios,
		XYRatio:   containerRatio,
		NumLeaves: numLeaves,
		Seed:      seed,
	}
}

func NewRandomBasic3DTreeGenerator(ratios TreeRatios, xyRatio, zyRatio float64, numLeaves int, seed int64) *RandomBasicTreeGenerator {
	return &RandomBasicTreeGenerator{
		Ratios:    ratios,
		XYRatio:   xyRatio,
		ZYRatio:   zyRatio,
		NumLeaves: numLeaves,
		Seed:      seed,
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
		"Ratios":               RatiosParameterString(gen.Ratios.Ratios()),
		"Container Ratio (XY)": strconv.FormatFloat(gen.XYRatio, 'f', 4, 64),
		"Container Ratio (ZY)": strconv.FormatFloat(gen.ZYRatio, 'f', 4, 64),
		"Number of Leaves":     gen.NumLeaves,
		"Random Seed":          gen.Seed,
	}
}

func (gen *RandomBasicTreeGenerator) Is3D() bool {
	return gen.ZYRatio > 0
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
	// We have horizontal and vertical splits defined in the complements array.
	// We have 3 possible planes that could be divided vertically/horizontally.
	xyRatioIndex := leaf.RatioIndexXY
	zyRatioIndex := leaf.RatioIndexZY

	xyComplements := complements[xyRatioIndex]
	zyComplements := complements[zyRatioIndex]

	// Visit each horizontal/vertical and see if it is compatible with the zy and
	// zx planes.

	// Any horizontal cut on the xy axis, affects the zy axis:
	//   XYCutHeight = XYRatio / XYRatioTop
	//   ZYRatioTop = ZYRatio / XYCutHeight
	//   Compatible if ZYRatioTop can be found in the index
	// Any vertical cut on the xy axis, affects the zx axis:
	//   XYCutWidth = XYRatioLeft / XYRatio
	//   XZRatioTop = XZRatio / XYCutWidth
	// Any vertical cut on the zy axis, affects the zx axis
	//   ZYCutWidth = ZYRatioLeft / ZYRatio
	//   XZRatioLeft = ZYCutWidth / XZRatio

	// Take each vertical split possible for the zyRatioIndex and check it against
	// the zx plane. If good, then add these to the possibilities as a DepthSplit
	// (instead of a vertical split)

	xyRatio := leaf.tree.Ratio(xyRatioIndex)
	zyRatio := leaf.tree.Ratio(zyRatioIndex)
	xzRatio := zyRatio / xyRatio
	var splits []Split

	for _, xySplit := range xyComplements {
		if xySplit.IsHorizontal() {
			cutHeight := RatioNormalHeight(xyRatio, leaf.tree.Ratio(xySplit.LeftIndex()))
			zyRatioTop := zyRatio / cutHeight
			index := FindClosestIndex(leaf.tree.ratios.Ratios(), zyRatioTop, leaf.tree.epsilon)
			if index < 0 {
				continue
			}
			// TODO: do we need to check the right as well? or is it guaranteed
		} else if xySplit.IsVertical() {
			cutWidth := RatioNormalWidth(xyRatio, leaf.tree.Ratio(xySplit.LeftIndex()))
			xzRatioTop := xzRatio / cutWidth
			index := FindClosestIndex(leaf.tree.ratios.Ratios(), xzRatioTop, leaf.tree.epsilon)
			if index < 0 {
				continue
			}
			// TODO: do we need to check the right as well?
		} else {
			panic("What type?")
		}

		splits = append(splits, xySplit)
	}

	for _, zySplit := range zyComplements {
		if !zySplit.IsVertical() {
			continue
		}

		cutWidth := RatioNormalWidth(zyRatio, leaf.tree.Ratio(zySplit.LeftIndex()))
		xzRatioLeft := cutWidth / xzRatio
		index := FindClosestIndex(leaf.tree.ratios.Ratios(), xzRatioLeft, leaf.tree.epsilon)
		if index < 0 {
			continue
		}
		// TODO: do we need to check the right as well?
		splits = append(splits, NewDepthSplit(zySplit.LeftIndex(), zySplit.RightIndex()))
	}

	if len(splits) == 0 {
		return nil
	}

	return &leafSplits{
		leaf:   leaf,
		splits: splits,
	}
}

func (gen *RandomBasicTreeGenerator) Generate() (*Tree, error) {
	rand.Seed(gen.Seed)

	epsilon := CalculateRatiosEpsilon(gen.Ratios.Ratios())
	xyRatioIndex := FindClosestIndex(gen.Ratios.Ratios(), gen.XYRatio, epsilon)
	if xyRatioIndex < 0 {
		return nil, errors.New("Container ratio not found in list of ratios.")
	}

	complements := gen.Ratios.Complements()

	// Generate the container
	var tree *Tree
	if !gen.Is3D() {
		tree = NewTree2D(gen.Ratios, xyRatioIndex)
	} else {
		zyRatioIndex := FindClosestIndex(gen.Ratios.Ratios(), gen.ZYRatio, epsilon)
		if zyRatioIndex < 0 {
			return nil, errors.New("Container ratio not found in list of ratios.")
		}
		tree = NewTree(gen.Ratios, xyRatioIndex, zyRatioIndex)
	}

	leafCount := 1

	leafDims := []*DimensionalNode{
		NewDimensionalNodeFromTree(tree, &Vector{0, 0, 0}, 1.0),
	}

	for {
		if leafCount >= gen.NumLeaves {
			break
		}

		it := NewDimensionalIteratorFromLeaves(leafDims)

		leafDims = leafDims[:0]
		for it.HasNext() {
			dimNode := it.Next()
			if dimNode.IsLeaf() {
				leafDims = append(leafDims, dimNode)
			}
		}

		// Collect all splittable leaves
		var filteredLeaves []*leafSplits

		for _, leaf := range leafDims {

			if !gen.Is3D() {
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
