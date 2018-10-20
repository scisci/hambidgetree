package hambidgetree

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("Not Found")

const UnityScale = 1.0

type Region struct {
	dimension    *Dimension
	ratioIndexXY int
	ratioIndexZY int
}

func NewRegion(dimension *Dimension, ratioIndexXY, ratioIndexZY int) *Region {
	return &Region{
		dimension:    dimension,
		ratioIndexXY: ratioIndexXY,
		ratioIndexZY: ratioIndexZY,
	}
}

func (region *Region) Dimension() *Dimension {
	return region.dimension
}

func (region *Region) RatioIndexXY() int {
	return region.ratioIndexXY
}

func (region *Region) RatioIndexZY() int {
	return region.ratioIndexZY
}

type NodeRegion interface {
	Region() *Region
	Node() Node
}

// Split the given region using the given ratioIndex
func SplitRegionHorizontal(ratios Ratios, region *Region, leftIndex, rightIndex int) (left, right *Region) {
	epsilon := 0.0000001

	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()

	leftRatio := ratios[leftIndex]
	leftHeightParam := RatioNormalHeight(ratios[ratioIndexXY], leftRatio)

	leftRatioIndexZY := ratioIndexZY
	rightRatioIndexZY := ratioIndexZY

	if IsRatioIndexDefined(ratioIndexZY) {
		zyRatio := ratios[ratioIndexZY]
		leftRatio := zyRatio / leftHeightParam
		rightRatio := zyRatio / (1 - leftHeightParam)
		leftRatioIndexZY = FindClosestIndexWithinRange(ratios, leftRatio, epsilon)
		rightRatioIndexZY = FindClosestIndexWithinRange(ratios, rightRatio, epsilon)
		if leftRatioIndexZY < 0 || rightRatioIndexZY < 0 {
			fmt.Printf("Failed to split horizontal container:%f, left:%f, right:%f\n", zyRatio, leftRatio, rightRatio)
			fmt.Printf("Got %d, %d in %v \n", leftRatioIndexZY, rightRatioIndexZY, ratios)
			panic("ZY Ratio is not one of the supported ratios!")
		}
	}

	// Find the right ratio index

	right = NewRegion(
		dimension.Inset(AxisY, dimension.Height()*leftHeightParam),
		rightIndex,
		rightRatioIndexZY,
	)

	left = NewRegion(
		dimension.Inset(AxisY, -dimension.Height()*(1-leftHeightParam)),
		leftIndex,
		leftRatioIndexZY,
	)

	return
}

func SplitRegionVertical(ratios Ratios, region *Region, leftIndex, rightIndex int) (left, right *Region) {
	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()

	leftRatio := ratios[leftIndex]
	leftWidthParam := RatioNormalWidth(ratios[ratioIndexXY], leftRatio) // left.Ratio() / ratio

	right = NewRegion(
		dimension.Inset(AxisX, dimension.Width()*leftWidthParam),
		rightIndex,
		ratioIndexZY,
	)

	left = NewRegion(
		dimension.Inset(AxisX, -dimension.Width()*(1-leftWidthParam)),
		leftIndex,
		ratioIndexZY,
	)

	return
}

func SplitRegionDepth(ratios Ratios, region *Region, leftIndex, rightIndex int) (left, right *Region) {
	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()
	//fmt.Println("splitting depth")
	leftRatio := ratios[leftIndex]
	leftDepthParam := RatioNormalWidth(ratios[ratioIndexZY], leftRatio)
	//fmt.Printf("depth: %f, container: %f, ratio: %f\n", dimension.Depth(), node.RatioZY, leftRatio)

	right = NewRegion(
		dimension.Inset(AxisZ, dimension.Depth()*leftDepthParam),
		ratioIndexXY,
		rightIndex,
	)

	left = NewRegion(
		dimension.Inset(AxisZ, -dimension.Depth()*(1-leftDepthParam)),
		ratioIndexXY,
		leftIndex,
	)

	return
}

type RegionIterator struct {
	tree    Tree
	regions []*nodeRatioRegion
}

type nodeRatioRegion struct {
	node   Node
	region *Region
}

func (region *nodeRatioRegion) Node() Node {
	return region.node
}

func (region *nodeRatioRegion) Region() *Region {
	return region.region
}

func NewRegionIterator(tree Tree, offset *Vector, scale float64) *RegionIterator {
	ratios := tree.RatioSource().Ratios()

	ratioIndexXY := tree.RatioIndexXY()
	ratioIndexZY := tree.RatioIndexZY()

	ratioXY := ratios[ratioIndexXY]
	ratioZY := 0.0

	if IsRatioIndexDefined(ratioIndexZY) {
		ratioZY = ratios[ratioIndexZY]
	}

	max := NewVector(ratioXY*scale, 1*scale, ratioZY*scale)

	root := tree.Root()

	region := &nodeRatioRegion{
		root,
		NewRegion(
			NewDimension3DV(offset, offset.Add(max)),
			ratioIndexXY,
			ratioIndexZY,
		),
	}

	return &RegionIterator{
		tree:    tree,
		regions: []*nodeRatioRegion{region},
	}
}

func (it *RegionIterator) HasNext() bool {
	return len(it.regions) > 0
}

func (it *RegionIterator) Next() NodeRegion {
	if !it.HasNext() {
		return nil
	}

	node := it.regions[len(it.regions)-1]
	it.regions = it.regions[:len(it.regions)-1]

	branch := node.node.Branch()
	ratios := it.tree.RatioSource().Ratios()

	if branch != nil {
		left := branch.Left()
		right := branch.Right()
		splitType := branch.SplitType()

		var leftRegion, rightRegion *Region
		if splitType == SplitTypeHorizontal {
			leftRegion, rightRegion = SplitRegionHorizontal(ratios,
				node.Region(),
				branch.LeftIndex(),
				branch.RightIndex())
		} else if splitType == SplitTypeVertical {
			// When we split vertically
			leftRegion, rightRegion = SplitRegionVertical(ratios,
				node.Region(),
				branch.LeftIndex(),
				branch.RightIndex())
		} else if splitType == SplitTypeDepth {
			leftRegion, rightRegion = SplitRegionDepth(ratios,
				node.Region(),
				branch.LeftIndex(),
				branch.RightIndex())
		} else {
			panic("Unknown split type")
		}

		it.regions = append(it.regions, &nodeRatioRegion{right, rightRegion})
		it.regions = append(it.regions, &nodeRatioRegion{left, leftRegion})
	}

	return node
}

func NewTreeRegionMap(tree Tree, offset *Vector, scale float64) RegionMap {
	lookup := make(map[NodeID]*Region)
	it := NewRegionIterator(tree, offset, scale)
	for it.HasNext() {
		region := it.Next()
		lookup[region.Node().ID()] = region.Region()
	}
	return lookup
}
