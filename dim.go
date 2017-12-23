package hambidgetree

type Dimension struct {
	x Extent
	y Extent
}

func NewDimension(left, top, right, bottom float64) *Dimension {
	return &Dimension{
		x: NewExtent(left, right),
		y: NewExtent(top, bottom),
	}
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

func (dim *Dimension) Width() float64 {
	return dim.Right() - dim.Left()
}

func (dim *Dimension) Height() float64 {
	return dim.Bottom() - dim.Top()
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

	return dim.y.Intersect(other.y)
}

// Calculates the amount the provided dimension/rect overlaps this rect on
// the left side. It returns 0 if there is no overlap. For its right side,
// should be within epsilon distance of this ones left side.
func (dim *Dimension) IntersectRight(other *Dimension, epsilon float64) Extent {
	return other.IntersectLeft(dim, epsilon)
}
