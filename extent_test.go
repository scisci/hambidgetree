package hambidgetree

import "testing"
import "math"

var extentIntersectionTests = []struct {
	e1           Extent
	e2           Extent
	intersection Extent
}{
	{ // Identity
		NewExtent(0, 1),
		NewExtent(0, 1),
		NewExtent(0, 1),
	},
	{ // Partial Overlap contained
		NewExtent(0, 1),
		NewExtent(0.5, 1),
		NewExtent(0.5, 1),
	},
	{ // Partial Overlap not contained
		NewExtent(0, 1),
		NewExtent(0.5, 2),
		NewExtent(0.5, 1),
	},
	{ // Partial Overlap Lower Bound Contained
		NewExtent(0, 1),
		NewExtent(0, 0.5),
		NewExtent(0, 0.5),
	},
	{ // Partial Overlap Lower Bound Not Contained
		NewExtent(0, 1),
		NewExtent(-1, 0.5),
		NewExtent(0, 0.5),
	},
	{ // Non Overlap Lower Bound
		NewExtent(0, 1),
		NewExtent(-1, -0.5),
		NewExtent(0, 0),
	},
	{ // Non Overlap Upper Bound
		NewExtent(0, 1),
		NewExtent(1.5, 2.0),
		NewExtent(1, 1),
	},
}

func TestExtentIntersection(t *testing.T) {
	for i, args := range extentIntersectionTests {
		intersection := args.e1.Intersect(args.e2)
		if !intersection.Equal(args.intersection) {
			t.Errorf("Intersection %d not equal expected %v got %v", i, args.intersection, intersection)
		}
	}
}

func TestExtentEmpty(t *testing.T) {
	empty := NewExtent(0.5, 0.5)
	notEmpty := NewExtent(-1, 1)
	if !empty.Empty() {
		t.Errorf("Extent should be empty, got %v", empty)
	}
	if notEmpty.Empty() {
		t.Errorf("Extent should not be empty, got %v", notEmpty)
	}
}

func TestExtentSize(t *testing.T) {
	e1 := NewExtent(-10, 10)
	if s := e1.Size(); s != 20.0 {
		t.Errorf("Extent should be 20, got %f", s)
	}

	e2 := NewExtent(10, 10)
	if s := e2.Size(); s != 0.0 {
		t.Errorf("Extent should be 0, got %f", s)
	}
}

func TestExtentEqual(t *testing.T) {
	e1 := NewExtent(-10, 10)
	e2 := NewExtent(-10, 10)
	e3 := NewExtent(-10, 0)
	e4 := NewExtent(0, 10)

	equal := func(lhs, rhs Extent) {
		if !lhs.Equal(rhs) {
			t.Errorf("Should be equal %v, %v", lhs, rhs)
		}
	}

	notEqual := func(lhs, rhs Extent) {
		if lhs.Equal(rhs) {
			t.Errorf("Should not be equal %v, %v", lhs, rhs)
		}
	}

	equal(e1, e2)
	equal(e2, e1)
	notEqual(e1, e3)
	notEqual(e3, e1)
	notEqual(e1, e4)
	notEqual(e4, e1)
	notEqual(e2, e3)
	notEqual(e3, e2)
	notEqual(e2, e4)
	notEqual(e4, e2)
}

var extentDistanceTests = []struct {
	e1       Extent
	e2       Extent
	distance float64
}{
	{ // Identity
		NewExtent(0, 1),
		NewExtent(0, 1),
		0,
	},
	{ // Overlapping
		NewExtent(0, 1),
		NewExtent(0.5, 1.5),
		0,
	},
	{ // Edge to Edge
		NewExtent(0, 1),
		NewExtent(1, 2),
		0,
	},
	{ // Separated
		NewExtent(0, 1),
		NewExtent(1.1, 1.2),
		0.1,
	},
	{ // Separated
		NewExtent(-1, -0.1),
		NewExtent(0.2, 0.4),
		0.3,
	},
}

func TestExtentDistance(t *testing.T) {
	for i, args := range extentDistanceTests {
		dist := args.e1.Distance(args.e2)
		dist2 := args.e2.Distance(args.e1)
		if math.Abs(args.distance-dist) > 0.0000001 {
			t.Errorf("Distance %d not equal expected %f got %f", i, args.distance, dist)
		}
		if math.Abs(dist2-dist) > 0.0000001 {
			t.Errorf("Distance %d inverse failed", i)
		}
	}
}
