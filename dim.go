package hambidgetree

import "fmt"

// Enum used to identify an axis to operate on for a given set of dimensions
type Axis int

const AxisX = 1
const AxisY = 2
const AxisZ = 4

type Vector struct {
	x float64
	y float64
	z float64
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{x, y, z}
}

func (v *Vector) Add(other *Vector) *Vector {
	return NewVector(v.x+other.x, v.y+other.y, v.z+other.z)
}

type Dimension struct {
	x Extent
	y Extent
	z Extent
}

func NewDimension2D(left, top, right, bottom float64) *Dimension {
	return &Dimension{
		x: NewExtent(left, right),
		y: NewExtent(top, bottom),
	}
}

func NewDimension3D(left, top, front, right, bottom, back float64) *Dimension {
	return &Dimension{
		x: NewExtent(left, right),
		y: NewExtent(top, bottom),
		z: NewExtent(front, back),
	}
}

func NewDimension3DV(min, max *Vector) *Dimension {
	return NewDimension3D(min.x, min.y, min.z, max.x, max.y, max.z)
}

func (dim *Dimension) Clone() *Dimension {
	return NewDimension3D(
		dim.x.start, dim.y.start, dim.z.start,
		dim.x.end, dim.y.end, dim.z.end)
}

func (dim *Dimension) Is3D() bool {
	return !dim.z.Empty()
}

func (dim *Dimension) String() string {
	return fmt.Sprintf("Dim{%.2f, %.2f, %.2f, %.2f}",
		dim.Left(), dim.Top(), dim.Width(), dim.Height())
}

func (dim *Dimension) Left() float64 {
	return dim.x.start
}

func (dim *Dimension) Right() float64 {
	return dim.x.end
}

func (dim *Dimension) Top() float64 {
	return dim.y.start
}

func (dim *Dimension) Bottom() float64 {
	return dim.y.end
}

func (dim *Dimension) Front() float64 {
	return dim.z.start
}

func (dim *Dimension) Back() float64 {
	return dim.z.end
}

func (dim *Dimension) Width() float64 {
	return dim.x.Size()
}

func (dim *Dimension) Height() float64 {
	return dim.y.Size()
}

func (dim *Dimension) Depth() float64 {
	return dim.z.Size()
}

// Insets the extent corresponding to the given axis. If its a positive value
// the start of the axis is inset, if its a negative value the end is inset.
func (dim *Dimension) Inset(axis Axis, distance float64) *Dimension {
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

func (dim *Dimension) DistanceSquared(other *Dimension) float64 {
	dx := dim.x.Distance(other.x)
	dy := dim.y.Distance(other.y)
	dz := dim.z.Distance(other.z)
	return dx*dx + dy*dy + dz*dz
}

// Calculates the amount the provided dimension/rect overlaps this rect on
// the left side. It returns 0 if there is no overlap. For its right side,
// should be within epsilon distance of this ones left side.
func (dim *Dimension) IntersectLeft(other *Dimension, epsilon float64) Extent {
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
func (dim *Dimension) IntersectRight(other *Dimension, epsilon float64) Extent {
	return other.IntersectLeft(dim, epsilon)
}
