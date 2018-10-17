package hambidgetree

import (
	"fmt"
	"math"
	"errors"
)

var ErrInvalidExtent = errors.New("Invalid extent")

// Extent
//
// Represents a range along a dimension. Start is included and end is not 
// included.
type Extent struct {
	start, end float64
}

func NewExtent(start, end float64) Extent {
	if (end < start) {
		panic(ErrInvalidExtent)
	}

	return Extent{
		start: start,
		end:   end,
	}
}

func (extent Extent) String() string {
	return fmt.Sprintf("Extent{%.2f, %.2f}", extent.start, extent.end)
}

func (extent Extent) Equal(other Extent) bool {
	return extent.start == other.start && extent.end == other.end
}

func (extent Extent) Empty() bool {
	return extent.end <= extent.start
}

func (extent Extent) NearlyEmpty(epsilon float64) bool {
	return extent.end - extent.start < epsilon
}

func (extent Extent) Size() float64 {
	return extent.end - extent.start
}

func (extent Extent) Distance(other Extent) float64 {
	if other.end < extent.start {
		return extent.start - other.end
	}

	if other.start > extent.end {
		return other.start - extent.end
	}

	return 0.0
}

func (extent Extent) Intersect(other Extent) Extent {
	if other.end < extent.start {
		return NewExtent(extent.start, extent.start)
	}

	if other.start > extent.end {
		return NewExtent(extent.end, extent.end)
	}

	return NewExtent(
		math.Max(extent.start, other.start),
		math.Min(extent.end, other.end))
}
