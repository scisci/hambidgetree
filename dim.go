package hambidgetree

import "fmt"

// Enum used to identify an axis to operate on for a given set of dimensions
type Axis int

const AxisX = 1
const AxisY = 2
const AxisZ = 3

type AlignedBox struct {
	x Extent
	y Extent
	z Extent
}

func NewAlignedBox2D(left, top, right, bottom float64) *AlignedBox {
	return &AlignedBox{
		x: NewExtent(left, right),
		y: NewExtent(top, bottom),
	}
}

func NewAlignedBox3D(left, top, front, right, bottom, back float64) *AlignedBox {
	return &AlignedBox{
		x: NewExtent(left, right),
		y: NewExtent(top, bottom),
		z: NewExtent(front, back),
	}
}

func NewAlignedBox3DV(min, max *Vector) *AlignedBox {
	return NewAlignedBox3D(min.x, min.y, min.z, max.x, max.y, max.z)
}

func (dim *AlignedBox) Clone() *AlignedBox {
	return NewAlignedBox3D(
		dim.x.start, dim.y.start, dim.z.start,
		dim.x.end, dim.y.end, dim.z.end)
}

func (dim *AlignedBox) Is3D() bool {
	return !dim.z.Empty()
}

func (dim *AlignedBox) String() string {
	return fmt.Sprintf("Dim{%.2f, %.2f, %.2f, %.2f}",
		dim.Left(), dim.Top(), dim.Width(), dim.Height())
}

func (dim *AlignedBox) Left() float64 {
	return dim.x.start
}

func (dim *AlignedBox) Right() float64 {
	return dim.x.end
}

func (dim *AlignedBox) Top() float64 {
	return dim.y.start
}

func (dim *AlignedBox) Bottom() float64 {
	return dim.y.end
}

func (dim *AlignedBox) Front() float64 {
	return dim.z.start
}

func (dim *AlignedBox) Back() float64 {
	return dim.z.end
}

func (dim *AlignedBox) Width() float64 {
	return dim.x.Size()
}

func (dim *AlignedBox) Height() float64 {
	return dim.y.Size()
}

func (dim *AlignedBox) Depth() float64 {
	return dim.z.Size()
}

// Insets the extent corresponding to the given axis. If its a positive value
// the start of the axis is inset, if its a negative value the end is inset.
func (dim *AlignedBox) Inset(axis Axis, distance float64) *AlignedBox {
	inset := dim.Clone()

	var extent *Extent
	switch axis {
	case AxisX:
		extent = &inset.x
	case AxisY:
		extent = &inset.y
	case AxisZ:
		extent = &inset.z
	}

	if distance > 0 {
		extent.start += distance
	} else {
		extent.end += distance
	}

	return inset
}

func (dim *AlignedBox) DistanceSquared(other *AlignedBox) float64 {
	dx := dim.x.Distance(other.x)
	dy := dim.y.Distance(other.y)
	dz := dim.z.Distance(other.z)
	return dx*dx + dy*dy + dz*dz
}

// Calculates the amount the provided dimension/rect overlaps this rect on
// the left side. It returns 0 if there is no overlap. For its right side,
// should be within epsilon distance of this ones left side.
func (dim *AlignedBox) IntersectLeft(other *AlignedBox, epsilon float64) Extent {
	leftDist := other.Right() - dim.Left()
	if leftDist > epsilon || leftDist < -epsilon {
		return Extent{
			start: dim.y.start,
			end:   dim.y.start,
		}
	}

	e := dim.y.Intersect(other.y)
	if e.NearlyEmpty(epsilon) {
		e.end = e.start
	}
	return e
}

// Calculates the amount the provided dimension/rect overlaps this rect on
// the left side. It returns 0 if there is no overlap. For its right side,
// should be within epsilon distance of this ones left side.
func (dim *AlignedBox) IntersectRight(other *AlignedBox, epsilon float64) Extent {
	return other.IntersectLeft(dim, epsilon)
}
