package randombasic

import (
	"errors"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/builder"
	"math/rand"
	"strconv"
)

const defaultEpsilon = 0.0000001

type RandomBasicTreeGenerator struct {
	NumLeaves   int
	RatioSource htree.RatioSource
	Complements htree.Complements
	Seed        int64
	XYRatio     float64
	ZYRatio     float64
}

type leafSplits struct {
	leaf   htree.Leaf
	splits []htree.Split
}

func IsSplitHorizontal(split htree.Split) bool {
	return split.Type() == htree.SplitTypeHorizontal
}

func IsSplitVertical(split htree.Split) bool {
	return split.Type() == htree.SplitTypeVertical
}

func IsSplitDepth(split htree.Split) bool {
	return split.Type() == htree.SplitTypeDepth
}

func New(ratioSource htree.RatioSource, containerRatio float64, numLeaves int, seed int64) (*RandomBasicTreeGenerator, error) {
	complements, err := htree.NewComplements(ratioSource.Ratios(), defaultEpsilon)
	if err != nil {
		return nil, err
	}

	return &RandomBasicTreeGenerator{
		NumLeaves:   numLeaves,
		RatioSource: ratioSource,
		Complements: complements,
		Seed:        seed,
		XYRatio:     containerRatio,
	}, nil
}

func New3D(ratioSource htree.RatioSource, xyRatio, zyRatio float64, numLeaves int, seed int64) (*RandomBasicTreeGenerator, error) {
	complements, err := htree.NewComplements(ratioSource.Ratios(), defaultEpsilon)
	if err != nil {
		return nil, err
	}

	return &RandomBasicTreeGenerator{
		NumLeaves:   numLeaves,
		RatioSource: ratioSource,
		Complements: complements,
		Seed:        seed,
		XYRatio:     xyRatio,
		ZYRatio:     zyRatio,
	}, nil
}

func (gen *RandomBasicTreeGenerator) Is3D() bool {
	return gen.ZYRatio > 0
}

func (gen *RandomBasicTreeGenerator) filterLeaves2D(leaf htree.Leaf, complements htree.Complements) *leafSplits {
	ratioIndex := leaf.RatioIndexXY() //   tree.RatioIndex(leaf.Node, htree.RatioPlaneXY)

	if len(complements[ratioIndex]) == 0 {
		return nil
	}

	return &leafSplits{
		leaf:   leaf,
		splits: complements[ratioIndex],
	}
}

func (gen *RandomBasicTreeGenerator) filterLeaves3D(leaf htree.Leaf, complements htree.Complements) *leafSplits {
	ratios := gen.RatioSource.Ratios()
	// We have horizontal and vertical splits defined in the complements array.
	// We have 3 possible planes that could be divided vertically/horizontally.
	xyRatioIndex := leaf.RatioIndexXY()
	zyRatioIndex := leaf.RatioIndexZY()

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

	xyRatio := ratios[xyRatioIndex]
	zyRatio := ratios[zyRatioIndex]
	zxRatio := zyRatio / xyRatio
	var splits []htree.Split

	for _, xySplit := range xyComplements {
		if IsSplitHorizontal(xySplit) {
			cutHeight := htree.RatioNormalHeight(xyRatio, ratios[xySplit.LeftIndex()])
			compHeight := htree.RatioNormalHeight(xyRatio, ratios[xySplit.RightIndex()])
			zyRatioTop := zyRatio / cutHeight
			zyRatioBottom := zyRatio / compHeight
			index := htree.FindClosestIndexWithinRange(ratios, zyRatioTop, defaultEpsilon)
			if index < 0 {
				continue
			}

			index = htree.FindClosestIndexWithinRange(ratios, zyRatioBottom, defaultEpsilon)
			if index < 0 {
				continue
				//fmt.Printf("tried to split h ratio %f, got %f and %f, but %f is not a valid ratio\n", xyRatio, zyRatioTop, zyRatioBottom, zyRatioBottom)
				//panic("right invalid")
			}
		} else if IsSplitVertical(xySplit) {
			cutWidth := htree.RatioNormalWidth(xyRatio, ratios[xySplit.LeftIndex()])
			compWidth := htree.RatioNormalWidth(xyRatio, ratios[xySplit.RightIndex()])
			zxRatioTop := zxRatio / cutWidth
			zxRatioBottom := zxRatio / compWidth
			index := htree.FindClosestIndexWithinRange(ratios, zxRatioTop, defaultEpsilon)
			if index < 0 {
				continue
			}

			index = htree.FindClosestIndexWithinRange(ratios, zxRatioBottom, defaultEpsilon)
			if index < 0 {
				continue
				//fmt.Printf("tried to split v ratio (xy:%f, xz:%f, zy:%f) with width %f from ratio %f against xz %f, got %f and %f, but %f is not a valid ratio\n", xyRatio, xzRatio, zyRatio, cutWidth, leaf.tree.Ratio(xySplit.LeftIndex()), xzRatio, xzRatioTop, xzRatioBottom, xzRatioBottom)
				//panic("right invalid")
			}
		} else {
			panic("What type?")
		}

		splits = append(splits, xySplit)
	}

	for _, zySplit := range zyComplements {
		if !IsSplitVertical(zySplit) {
			continue
		}

		cutWidth := htree.RatioNormalWidth(zyRatio, ratios[zySplit.LeftIndex()])
		compWidth := htree.RatioNormalWidth(zyRatio, ratios[zySplit.RightIndex()])
		zxRatioLeft := cutWidth * zxRatio
		index := htree.FindClosestIndexWithinRange(ratios, zxRatioLeft, defaultEpsilon)
		if index < 0 {
			continue
		}

		zxRatioRight := compWidth * zxRatio
		index = htree.FindClosestIndexWithinRange(ratios, zxRatioRight, defaultEpsilon)
		if index < 0 {
			continue
			//fmt.Printf("tried to split d ratio %f, got %f and %f, but %f is not a valid ratio\n", xyRatio, xzRatioLeft, xzRatioRight, xzRatioRight)
			//panic("right invalid")
		}
		splits = append(splits, htree.NewDepthSplit(zySplit.LeftIndex(), zySplit.RightIndex()))
	}

	if len(splits) == 0 {
		return nil
	}

	return &leafSplits{
		leaf:   leaf,
		splits: splits,
	}
}

func (gen *RandomBasicTreeGenerator) Generate() (htree.Tree, error) {
	rand.Seed(gen.Seed)

	ratios := gen.RatioSource.Ratios()

	epsilon := htree.CalculateRatiosEpsilon(ratios)
	xyRatioIndex := htree.FindClosestIndex(ratios, gen.XYRatio, epsilon)
	if xyRatioIndex < 0 {
		return nil, errors.New("Container ratio not found in list of ratios.")
	}

	complements, err := htree.NewComplements(ratios, defaultEpsilon)
	if err != nil {
		return nil, err
	}

	// Generate the container
	var treeBuilder *builder.TreeBuilder
	if !gen.Is3D() {
		treeBuilder = builder.New2D(gen.RatioSource, xyRatioIndex)
	} else {
		zyRatioIndex := htree.FindClosestIndex(ratios, gen.ZYRatio, epsilon)
		if zyRatioIndex < 0 {
			return nil, errors.New("Container ratio not found in list of ratios.")
		}
		treeBuilder = builder.New3D(gen.RatioSource, xyRatioIndex, zyRatioIndex)
	}

	//leafCount := 1

	//leafDims := treeBuilder.Leaves()

	for {
		leaves := treeBuilder.Leaves()

		if len(leaves) >= gen.NumLeaves {
			break
		}

		//it := htree.NewDimensionalIterator(tree, htree.Origin, htree.UnityScale) // NewDimensionalIteratorFromLeaves(leafDims)

		/*
			leafDims = leafDims[:0]
			for it.HasNext() {
				dimNode := it.Next()
				if dimNode.IsLeaf() {
					leafDims = append(leafDims, dimNode)
				}
			}
		*/

		// Collect all splittable leaves
		var filteredLeaves []*leafSplits

		for _, leaf := range leaves {

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
			return nil, errors.New("Unable to reach desired number of leaves (" +
				strconv.Itoa(gen.NumLeaves) + "), got " + strconv.Itoa(len(leaves)) + ".")
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
			split = htree.NewInvertedSplit(split)
		}

		treeBuilder.Branch(filteredLeaf.leaf.ID(), split.Type(), split.LeftIndex(), split.RightIndex())
		//filteredLeaf.leaf.Divide(split)
		//leafCount += 1
	}

	tree, _ := treeBuilder.Build()
	return tree, nil
}
