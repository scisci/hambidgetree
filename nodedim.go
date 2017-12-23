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
}

type DimensionalIterator struct {
	dimensions []*DimensionalNode
}

func NewDimensionalIterator(root *Node, offsetX, offsetY, scale float64) *DimensionalIterator {
	return &DimensionalIterator{
		dimensions: []*DimensionalNode{
			&DimensionalNode{
				root,
				NewDimension(offsetX, offsetY, root.Ratio()*scale, 1*scale),
			},
		},
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
		ratio := node.Ratio()
		left := node.Node.left
		right := node.Node.right

		if node.Split().IsHorizontal() {
			leftHeight := dimension.Height() * ratio / left.Ratio()

			it.dimensions = append(it.dimensions, &DimensionalNode{
				right,
				NewDimension(
					dimension.Left(),
					dimension.Top()+leftHeight,
					dimension.Right(),
					dimension.Bottom(),
				),
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				left,
				NewDimension(
					dimension.Left(),
					dimension.Top(),
					dimension.Right(),
					dimension.Top()+leftHeight,
				),
			})
		} else {
			leftWidth := dimension.Width() * left.Ratio() / ratio

			it.dimensions = append(it.dimensions, &DimensionalNode{
				right,
				NewDimension(
					dimension.Left()+leftWidth,
					dimension.Top(),
					dimension.Right(),
					dimension.Bottom(),
				),
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				left,
				NewDimension(
					dimension.Left(),
					dimension.Top(),
					dimension.Left()+leftWidth,
					dimension.Bottom(),
				),
			})
		}

	}

	return node
}

type NodeDimensionMap struct {
	lookup map[NodeID]*Dimension
}

func NewNodeDimensionMap(root *Node, offsetX, offsetY, scale float64) *NodeDimensionMap {
	lookup := make(map[NodeID]*Dimension)
	it := NewDimensionalIterator(root, offsetX, offsetY, scale)
	for it.HasNext() {
		dimNode := it.Next()
		lookup[dimNode.Node.id] = dimNode.Dimension
	}
	return &NodeDimensionMap{
		lookup: lookup,
	}
}

func (nodeDimMap NodeDimensionMap) Dimension(node *Node) (*Dimension, error) {
	if dim, ok := nodeDimMap.lookup[node.id]; !ok {
		return nil, ErrNotFound
	} else {
		return dim, nil
	}
}
