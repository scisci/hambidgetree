package hambidgetree

import (
	"errors"
)

var ErrNotFound = errors.New("Not Found")

type NodeDimensions interface {
	Dimension(node *Node) (*Dimension, error)
}

type DimensionalNode struct {
	*Node
	*Dimension
	RatioIndexXY int
	RatioIndexZY int
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
			leftRatioIndexZY := FindClosestIndex(node.tree.ratios.Ratios(), node.tree.Ratio(node.RatioIndexZY)/leftHeightParam, node.tree.epsilon)
			rightRatioIndexZY := FindClosestIndex(node.tree.ratios.Ratios(), node.tree.Ratio(node.RatioIndexZY)/(1-leftHeightParam), node.tree.epsilon)
			if leftRatioIndexZY < 0 || rightRatioIndexZY < 0 {
				panic("ZY Ratio is not one of the supported ratios!")
			}
			// When we split horizontally we modify both the xy plane AND the yz plane
			leftHeight := dimension.Height() * leftHeightParam // ratio / left.Ratio()

			//fmt.Printf("split horizontal container: %f ratio: %f\n", node.RatioZY, leftRatioXY)
			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         right,
				Dimension:    dimension.Inset(AxisY, leftHeight),
				RatioIndexXY: right.RatioIndex(),
				RatioIndexZY: rightRatioIndexZY,
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         left,
				Dimension:    dimension.Inset(AxisY, -leftHeight),
				RatioIndexXY: left.RatioIndex(),
				RatioIndexZY: leftRatioIndexZY,
			})
		} else if split.IsVertical() {
			// When we split vertically
			leftWidth := dimension.Width() * RatioNormalWidth(node.tree.Ratio(node.RatioIndexXY), left.Ratio()) // left.Ratio() / ratio

			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         right,
				Dimension:    dimension.Inset(AxisX, leftWidth),
				RatioIndexXY: right.RatioIndex(),
				RatioIndexZY: node.RatioIndexZY,
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         left,
				Dimension:    dimension.Inset(AxisX, -leftWidth),
				RatioIndexXY: left.RatioIndex(),
				RatioIndexZY: node.RatioIndexZY,
			})
		} else if split.IsDepth() {
			//fmt.Println("splitting depth")
			leftDepth := dimension.Depth() * RatioNormalWidth(node.tree.Ratio(node.RatioIndexZY), left.Ratio())
			//fmt.Printf("depth: %f, container: %f, ratio: %f\n", dimension.Depth(), node.RatioZY, leftRatio)

			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         right,
				Dimension:    dimension.Inset(AxisZ, leftDepth),
				RatioIndexXY: node.RatioIndexXY,
				RatioIndexZY: right.RatioIndex(),
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				Node:         left,
				Dimension:    dimension.Inset(AxisZ, -leftDepth),
				RatioIndexXY: node.RatioIndexXY,
				RatioIndexZY: left.RatioIndex(),
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

func (nodeDimMap NodeDimensionMap) Dimension(node *Node) (*Dimension, error) {
	if dim, ok := nodeDimMap.lookup[node.id]; !ok {
		return nil, ErrNotFound
	} else {
		return dim.Dimension, nil
	}
}
