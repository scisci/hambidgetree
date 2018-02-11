package hambidgetree

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("Not Found")

var Origin = &Vector{0.0, 0.0, 0.0}

const UnityScale = 1.0

type NodeDimensions interface {
	Dimension(id NodeID) (*Dimension, error)
}

type DimensionalNode struct {
	*Node
	*Dimension
	ParentDimensionalNode *DimensionalNode
	RatioIndexXY          int
	RatioIndexZY          int
}

func NewDimensionalNodeFromTree(tree *Tree, offset *Vector, scale float64) *DimensionalNode {
	ratioXY := tree.RatioXY()
	ratioZY := tree.RatioZY()
	return &DimensionalNode{
		Node: tree.root,
		Dimension: NewDimension3D(
			offset.x,
			offset.y,
			offset.z,
			offset.x+ratioXY*scale,
			offset.y+1*scale,
			offset.z+ratioZY*scale),
		RatioIndexXY: tree.xyRatioIndex,
		RatioIndexZY: tree.zyRatioIndex,
	}
}

type DimensionalIterator struct {
	dimensions []*DimensionalNode
}

func NewDimensionalIterator(tree *Tree, offset *Vector, scale float64) *DimensionalIterator {
	return &DimensionalIterator{
		dimensions: []*DimensionalNode{
			NewDimensionalNodeFromTree(tree, offset, scale),
		},
	}
}

func NewDimensionalIteratorFromLeaves(leaves []*DimensionalNode) *DimensionalIterator {
	return &DimensionalIterator{
		dimensions: leaves,
	}
}

func (it *DimensionalIterator) HasNext() bool {
	return len(it.dimensions) > 0
}

type Region interface {
	Dimension() *Dimension
	RatioIndexXY() int
	RatioIndexZY() int
}

type NodeRegion interface {
	Region
	Node() ImmutableNode
}

type SimpleRatioRegion struct {
	dimension    *Dimension
	ratioIndexXY int
	ratioIndexZY int
}

func (region *SimpleRatioRegion) Dimension() *Dimension {
	return region.dimension
}

func (region *SimpleRatioRegion) RatioIndexXY() int {
	return region.ratioIndexXY
}

func (region *SimpleRatioRegion) RatioIndexZY() int {
	return region.ratioIndexZY
}

func IsRatioIndexDefined(index int) bool {
	return index > RatioIndexUndefined
}

// Split the given region using the given ratioIndex
func SplitRegionHorizontal(ratios Ratios, region Region, leftIndex, rightIndex int) (left, right *SimpleRatioRegion) {
	epsilon := 0.0000001

	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()

	leftRatio := ratios.At(leftIndex)
	leftHeightParam := RatioNormalHeight(ratios.At(ratioIndexXY), leftRatio)

	leftRatioIndexZY := ratioIndexZY
	rightRatioIndexZY := ratioIndexZY

	if IsRatioIndexDefined(ratioIndexZY) {
		zyRatio := ratios.At(ratioIndexZY)
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

	right = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisY, dimension.Height()*leftHeightParam),
		ratioIndexXY: rightIndex,
		ratioIndexZY: rightRatioIndexZY,
	}

	left = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisY, -dimension.Height()*(1-leftHeightParam)),
		ratioIndexXY: leftIndex,
		ratioIndexZY: leftRatioIndexZY,
	}

	return
}

func SplitRegionVertical(ratios Ratios, region Region, leftIndex, rightIndex int) (left, right *SimpleRatioRegion) {
	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()

	leftRatio := ratios.At(leftIndex)
	leftWidthParam := RatioNormalWidth(ratios.At(ratioIndexXY), leftRatio) // left.Ratio() / ratio

	right = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisX, dimension.Width()*leftWidthParam),
		ratioIndexXY: rightIndex,
		ratioIndexZY: ratioIndexZY,
	}

	left = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisX, -dimension.Width()*(1-leftWidthParam)),
		ratioIndexXY: leftIndex,
		ratioIndexZY: ratioIndexZY,
	}

	return
}

func SplitRegionDepth(ratios Ratios, region Region, leftIndex, rightIndex int) (left, right *SimpleRatioRegion) {
	dimension := region.Dimension()
	ratioIndexXY := region.RatioIndexXY()
	ratioIndexZY := region.RatioIndexZY()
	//fmt.Println("splitting depth")
	leftRatio := ratios.At(leftIndex)
	leftDepthParam := RatioNormalWidth(ratios.At(ratioIndexZY), leftRatio)
	//fmt.Printf("depth: %f, container: %f, ratio: %f\n", dimension.Depth(), node.RatioZY, leftRatio)

	right = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisZ, dimension.Depth()*leftDepthParam),
		ratioIndexXY: ratioIndexXY,
		ratioIndexZY: rightIndex,
	}

	left = &SimpleRatioRegion{
		dimension:    dimension.Inset(AxisZ, -dimension.Depth()*(1-leftDepthParam)),
		ratioIndexXY: ratioIndexXY,
		ratioIndexZY: leftIndex,
	}

	return
}

func NewDimensionalNodeFromRegion(node *Node, region Region) *DimensionalNode {
	return &DimensionalNode{
		Node:         node,
		Dimension:    region.Dimension(),
		RatioIndexXY: region.RatioIndexXY(),
		RatioIndexZY: region.RatioIndexZY(),
	}
}

func (it *DimensionalIterator) Next() *DimensionalNode {
	if !it.HasNext() {
		return nil
	}

	node := it.dimensions[len(it.dimensions)-1]
	it.dimensions = it.dimensions[:len(it.dimensions)-1]

	dimension := node.Dimension

	if !node.IsLeaf() {
		left := node.Node.left
		right := node.Node.right
		split := node.Split()

		if split.IsHorizontal() {
			leftRegion, rightRegion := SplitRegionHorizontal(node.tree.ratios.Ratios(),
				&SimpleRatioRegion{
					dimension:    dimension,
					ratioIndexXY: node.RatioIndexXY,
					ratioIndexZY: node.RatioIndexZY,
				},
				left.RatioIndex(),
				right.RatioIndex())

			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(right, rightRegion))
			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(left, leftRegion))
			/*
				leftHeightParam := RatioNormalHeight(node.tree.Ratio(node.RatioIndexXY), left.Ratio())

				leftRatioIndexZY := node.RatioIndexZY
				rightRatioIndexZY := node.RatioIndexZY

				if node.RatioIndexZY > -1 {
					zyRatio := node.tree.Ratio(node.RatioIndexZY)
					leftRatio := zyRatio / leftHeightParam
					rightRatio := zyRatio / (1 - leftHeightParam)
					leftRatioIndexZY = FindClosestIndexWithinRange(node.tree.ratios.Ratios(), leftRatio, 0.0000001)
					rightRatioIndexZY = FindClosestIndexWithinRange(node.tree.ratios.Ratios(), rightRatio, 0.0000001)
					if leftRatioIndexZY < 0 || rightRatioIndexZY < 0 {
						fmt.Printf("Failed to split horizontal container:%f, left:%f, right:%f\n", zyRatio, leftRatio, rightRatio)
						panic("ZY Ratio is not one of the supported ratios!")
					}
				}

				//fmt.Printf("split horizontal container: %f ratio: %f\n", node.RatioZY, leftRatioXY)
				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  right,
					Dimension:             dimension.Inset(AxisY, dimension.Height()*leftHeightParam),
					ParentDimensionalNode: node,
					RatioIndexXY:          right.RatioIndex(),
					RatioIndexZY:          rightRatioIndexZY,
				})

				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  left,
					Dimension:             dimension.Inset(AxisY, -dimension.Height()*(1-leftHeightParam)),
					ParentDimensionalNode: node,
					RatioIndexXY:          left.RatioIndex(),
					RatioIndexZY:          leftRatioIndexZY,
				})
			*/
		} else if split.IsVertical() {
			// When we split vertically
			leftRegion, rightRegion := SplitRegionVertical(node.tree.ratios.Ratios(),
				&SimpleRatioRegion{
					dimension:    dimension,
					ratioIndexXY: node.RatioIndexXY,
					ratioIndexZY: node.RatioIndexZY,
				},
				left.RatioIndex(),
				right.RatioIndex())

			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(right, rightRegion))
			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(left, leftRegion))
			/*
				leftWidthParam := RatioNormalWidth(node.tree.Ratio(node.RatioIndexXY), left.Ratio()) // left.Ratio() / ratio

				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  right,
					Dimension:             dimension.Inset(AxisX, dimension.Width()*leftWidthParam),
					ParentDimensionalNode: node,
					RatioIndexXY:          right.RatioIndex(),
					RatioIndexZY:          node.RatioIndexZY,
				})

				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  left,
					Dimension:             dimension.Inset(AxisX, -dimension.Width()*(1-leftWidthParam)),
					ParentDimensionalNode: node,
					RatioIndexXY:          left.RatioIndex(),
					RatioIndexZY:          node.RatioIndexZY,
				})*/
		} else if split.IsDepth() {
			leftRegion, rightRegion := SplitRegionDepth(node.tree.ratios.Ratios(),
				&SimpleRatioRegion{
					dimension:    dimension,
					ratioIndexXY: node.RatioIndexXY,
					ratioIndexZY: node.RatioIndexZY,
				},
				left.RatioIndex(),
				right.RatioIndex())

			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(right, rightRegion))
			it.dimensions = append(it.dimensions, NewDimensionalNodeFromRegion(left, leftRegion))
			//fmt.Println("splitting depth")
			/*
				leftDepthParam := RatioNormalWidth(node.tree.Ratio(node.RatioIndexZY), left.Ratio())
				//fmt.Printf("depth: %f, container: %f, ratio: %f\n", dimension.Depth(), node.RatioZY, leftRatio)

				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  right,
					Dimension:             dimension.Inset(AxisZ, dimension.Depth()*leftDepthParam),
					ParentDimensionalNode: node,
					RatioIndexXY:          node.RatioIndexXY,
					RatioIndexZY:          right.RatioIndex(),
				})

				it.dimensions = append(it.dimensions, &DimensionalNode{
					Node:                  left,
					Dimension:             dimension.Inset(AxisZ, -dimension.Depth()*(1-leftDepthParam)),
					ParentDimensionalNode: node,
					RatioIndexXY:          node.RatioIndexXY,
					RatioIndexZY:          left.RatioIndex(),
				})
			*/
		}

	}

	return node
}

type RegionIterator struct {
	tree    ImmutableTree
	regions []*nodeRatioRegion
}

type nodeRatioRegion struct {
	ImmutableNode
	*SimpleRatioRegion
}

func (region *nodeRatioRegion) Node() ImmutableNode {
	return region.ImmutableNode
}

func NewRegionIterator(tree ImmutableTree, offset *Vector, scale float64) *RegionIterator {
	ratios := tree.Ratios()

	ratioIndexXY := tree.RatioIndexXY()
	ratioIndexZY := tree.RatioIndexZY()

	ratioXY := ratios.At(ratioIndexXY)
	ratioZY := 0.0

	if IsRatioIndexDefined(ratioIndexZY) {
		ratioZY = ratios.At(ratioIndexZY)
	}

	max := NewVector(ratioXY*scale, 1*scale, ratioZY*scale)

	root := tree.Root()

	region := &nodeRatioRegion{
		root,
		&SimpleRatioRegion{
			dimension:    NewDimension3DV(offset, offset.Add(max)),
			ratioIndexXY: ratioIndexXY,
			ratioIndexZY: ratioIndexZY,
		},
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

	branch := node.Branch()

	if branch != nil {
		left := branch.Left()
		right := branch.Right()
		splitType := branch.SplitType()

		var leftRegion, rightRegion *SimpleRatioRegion
		if splitType == SplitTypeHorizontal {
			leftRegion, rightRegion = SplitRegionHorizontal(it.tree.Ratios(),
				node.SimpleRatioRegion,
				branch.LeftIndex(),
				branch.RightIndex())
		} else if splitType == SplitTypeVertical {
			// When we split vertically
			leftRegion, rightRegion = SplitRegionVertical(it.tree.Ratios(),
				node.SimpleRatioRegion,
				branch.LeftIndex(),
				branch.RightIndex())
		} else if splitType == SplitTypeDepth {
			leftRegion, rightRegion = SplitRegionDepth(it.tree.Ratios(),
				node.SimpleRatioRegion,
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

type NodeRegionMap struct {
	offset *Vector
	scale  float64
	lookup map[NodeID]Region
}

func NewNodeRegionMap(tree ImmutableTree, offset *Vector, scale float64) *NodeRegionMap {
	lookup := make(map[NodeID]Region)
	it := NewRegionIterator(tree, offset, scale)
	for it.HasNext() {
		region := it.Next()
		lookup[region.Node().ID()] = region
	}
	return &NodeRegionMap{
		offset: offset,
		scale:  scale,
		lookup: lookup,
	}
}

func (m *NodeRegionMap) Offset() *Vector {
	return m.offset
}

func (m *NodeRegionMap) Scale() float64 {
	return m.scale
}

func (m *NodeRegionMap) Region(id NodeID) Region {
	return m.lookup[id]
}

type NodeDimensionMap struct {
	lookup map[NodeID]*DimensionalNode
}

func NewNodeDimensionMap(tree *Tree, offset *Vector, scale float64) *NodeDimensionMap {
	lookup := make(map[NodeID]*DimensionalNode)
	it := NewDimensionalIterator(tree, offset, scale)
	for it.HasNext() {
		dimNode := it.Next()
		lookup[dimNode.Node.id] = dimNode
	}
	return &NodeDimensionMap{
		lookup: lookup,
	}
}

func (nodeDimMap NodeDimensionMap) Dimension(id NodeID) (*Dimension, error) {
	if dim, ok := nodeDimMap.lookup[id]; !ok {
		return nil, ErrNotFound
	} else {
		return dim.Dimension, nil
	}
}
