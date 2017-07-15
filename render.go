package hambidgetree

type GraphicsContext interface {
	MoveTo(x, y float64)
	LineTo(x, y float64)
	Rect(x, y, w, h float64)
}

type Rect struct {
	x, y          float64
	width, height float64
}

type Dimension struct {
	left, top, right, bottom float64
}

func (dim Dimension) Width() float64 {
	return dim.right - dim.left
}

func (dim Dimension) Height() float64 {
	return dim.bottom - dim.top
}

type DimensionalNode struct {
	*HambidgeTreeNode
	Dimension
}

type DimensionalIterator struct {
	dimensions []*DimensionalNode
}

func NewDimensionalIterator(root *HambidgeTreeNode) *DimensionalIterator {
	return &DimensionalIterator{
		dimensions: []*DimensionalNode{
			&DimensionalNode{
				root,
				Dimension{0, 0, root.Ratio(), 1},
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
		left := node.HambidgeTreeNode.left
		right := node.HambidgeTreeNode.right

		if node.Split().IsHorizontal() {
			leftHeight := dimension.Height() * ratio / left.Ratio()

			it.dimensions = append(it.dimensions, &DimensionalNode{
				right,
				Dimension{
					dimension.left,
					dimension.top + leftHeight,
					dimension.right,
					dimension.bottom,
				},
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				left,
				Dimension{
					dimension.left,
					dimension.top,
					dimension.right,
					dimension.top + leftHeight,
				},
			})
		} else {
			leftWidth := dimension.Width() * left.Ratio() / ratio

			it.dimensions = append(it.dimensions, &DimensionalNode{
				right,
				Dimension{
					dimension.left + leftWidth,
					dimension.top,
					dimension.right,
					dimension.bottom,
				},
			})

			it.dimensions = append(it.dimensions, &DimensionalNode{
				left,
				Dimension{
					dimension.left,
					dimension.top,
					dimension.left + leftWidth,
					dimension.bottom,
				},
			})
		}

	}

	return node
}

type TreeStrokeRenderer struct {
}

func (renderer *TreeStrokeRenderer) Render(tree *HambidgeTree, gc GraphicsContext) {
	it := NewDimensionalIterator(tree.root)

	var container *DimensionalNode

	for it.HasNext() {
		node := it.Next()

		if container == nil {
			container = node
		}
		// Draw the stroke
		if !node.IsLeaf() {
			if node.Split().IsHorizontal() {
				y := node.Dimension.top + node.Dimension.Height()*node.HambidgeTreeNode.left.Ratio()
				gc.MoveTo(node.Dimension.left, y)
				gc.LineTo(node.Dimension.right, y)
			} else {
				x := node.Dimension.left + node.Dimension.Width()*node.HambidgeTreeNode.left.Ratio()
				gc.MoveTo(x, node.Dimension.top)
				gc.LineTo(x, node.Dimension.bottom)
			}
		}
	}

	// Finally draw the rectangle
	gc.Rect(0, 0, container.Width(), container.Height())

}
