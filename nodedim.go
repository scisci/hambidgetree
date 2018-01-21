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
		} else if split.IsVertical() {
			// When we split vertically
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
			})
		} else if split.IsDepth() {
			//fmt.Println("splitting depth")
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
		}

	}

	return node
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
