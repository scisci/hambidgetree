package hambidgetree

import "math"
import "strconv"

const RenderEpsilon = 0.0000001

type GraphicsContextCall interface {
	Name() string
	String() string
	Equals(other GraphicsContextCall) bool
}

type GraphicsContextLine struct {
	x1, y1, x2, y2 float64
}

func (m *GraphicsContextLine) Name() string {
	return "Line"
}

func (m *GraphicsContextLine) Equals(other GraphicsContextCall) bool {
	if l, ok := other.(*GraphicsContextLine); ok {
		return math.Abs(l.x1-m.x1) < RenderEpsilon && math.Abs(l.y1-m.y1) < RenderEpsilon && math.Abs(l.x2-m.x2) < RenderEpsilon && math.Abs(l.y2-m.y2) < RenderEpsilon
	}

	return false
}

func (m *GraphicsContextLine) String() string {
	return "Line{" +
		strconv.FormatFloat(m.x1, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.y1, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.x2, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.y2, 'f', -1, 64) + "}"
}

type GraphicsContextRect struct {
	x, y, width, height float64
}

func (m *GraphicsContextRect) Name() string {
	return "Rect"
}

func (m *GraphicsContextRect) String() string {
	return "Rect{" +
		strconv.FormatFloat(m.x, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.y, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.width, 'f', -1, 64) + ", " +
		strconv.FormatFloat(m.height, 'f', -1, 64) + "}"
}

func (m *GraphicsContextRect) Equals(other GraphicsContextCall) bool {
	if r, ok := other.(*GraphicsContextRect); ok {
		return math.Abs(r.x-m.x) < RenderEpsilon && math.Abs(r.y-m.y) < RenderEpsilon && math.Abs(r.width-m.width) < RenderEpsilon && math.Abs(r.height-m.height) < RenderEpsilon
	}

	return false
}

type GraphicsContext interface {
	Line(x1, y1, x2, y2 float64)
	Rect(x, y, w, h float64)
}

type GraphicsContextRecorder struct {
	Calls []GraphicsContextCall
}

func NewGraphicsContextRecorder() *GraphicsContextRecorder {
	return &GraphicsContextRecorder{}
}

func (gc *GraphicsContextRecorder) Line(x1, y1, x2, y2 float64) {
	gc.AddCall(&GraphicsContextLine{x1, y1, x2, y2})
}

func (gc *GraphicsContextRecorder) Rect(x, y, width, height float64) {
	gc.AddCall(&GraphicsContextRect{x, y, width, height})
}

func (gc *GraphicsContextRecorder) AddCall(call GraphicsContextCall) {
	gc.Calls = append(gc.Calls, call)
}

type TreeStrokeRenderer struct {
	offsetX float64
	offsetY float64
	scale   float64
	snap    bool
}

func NewTreeStrokeRenderer(offsetX, offsetY, scale float64) *TreeStrokeRenderer {
	return &TreeStrokeRenderer{
		offsetX: offsetX,
		offsetY: offsetY,
		scale:   scale,
		snap:    false,
	}
}

func (renderer *TreeStrokeRenderer) Snap(snap bool) {
	renderer.snap = snap
}

func (renderer *TreeStrokeRenderer) Render(tree *Tree, gc GraphicsContext) {
	it := NewDimensionalIterator(tree.root, renderer.offsetX, renderer.offsetY, renderer.scale)

	var container *DimensionalNode

	for it.HasNext() {
		node := it.Next()

		if container == nil {
			container = node
		}
		// Draw the stroke
		if !node.IsLeaf() {
			if node.Split().IsHorizontal() {
				y := node.Dimension.top + node.Dimension.Height()*node.Ratio()/node.Node.left.Ratio()
				if renderer.snap {
					y = math.Floor(y + 0.5)
				}

				gc.Line(node.Dimension.left, y, node.Dimension.right, y)
			} else {
				x := node.Dimension.left + node.Dimension.Width()*node.Node.left.Ratio()/node.Ratio()
				if renderer.snap {
					x = math.Floor(x + 0.5)
				}
				gc.Line(x, node.Dimension.top, x, node.Dimension.bottom)
			}
		}
	}

	// Finally draw the rectangle
	gc.Rect(container.Dimension.left, container.Dimension.top, container.Width(), container.Height())

}
