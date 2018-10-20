package hambidgetree

import "testing"
import "math"

var intersectionTests = []struct {
	d1                *AlignedBox
	d2                *AlignedBox
	leftIntersection  Extent
	rightIntersection Extent
}{
	{ // First block to the left of second block, same height
		NewAlignedBox2D(0, 0, 100, 100),
		NewAlignedBox2D(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 100),
	},
	{ // First block to the right of second block, same height
		NewAlignedBox2D(100, 0, 200, 100),
		NewAlignedBox2D(0, 0, 100, 100),
		NewExtent(0, 100),
		NewExtent(0, 0),
	},
	{ // First block to the left and shifted up from second block
		NewAlignedBox2D(0, -50, 100, 50),
		NewAlignedBox2D(100, 0, 200, 100),
		NewExtent(-50, -50),
		NewExtent(0, 50),
	},
	{ // First block to the left and shifted down from second block
		NewAlignedBox2D(0, 50, 100, 150),
		NewAlignedBox2D(100, 0, 200, 100),
		NewExtent(50, 50),
		NewExtent(50, 100),
	},
	{ // First block to the left but further than epsilon
		NewAlignedBox2D(0, 0, 50, 100),
		NewAlignedBox2D(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 0),
	},
	{ // First block to the right of second block but further than epsilon, same height
		NewAlignedBox2D(110, 0, 200, 100),
		NewAlignedBox2D(0, 0, 100, 100),
		NewExtent(0, 0),
		NewExtent(0, 0),
	},
}

func TestIntersect(t *testing.T) {
	for i, args := range intersectionTests {
		leftIntersection := args.d1.IntersectLeft(args.d2, 0.0000001)
		rightIntersection := args.d1.IntersectRight(args.d2, 0.0000001)
		if !leftIntersection.Equal(args.leftIntersection) {
			t.Errorf("Left intersection %d wrong, expected %v got %v", i, args.leftIntersection, leftIntersection)
		}
		if !rightIntersection.Equal(args.rightIntersection) {
			t.Errorf("Right intersection %d wrong, expected %v got %v", i, args.rightIntersection, rightIntersection)
		}
	}
}

func TestInset(t *testing.T) {
	dim := NewAlignedBox3D(0, 0, 0, 20, 20, 20)
	dim1 := dim.Inset(AxisX, 5)
	dim2 := dim.Inset(AxisX, -5)
	dim3 := dim.Inset(AxisY, 5)
	dim4 := dim.Inset(AxisY, -5)
	dim5 := dim.Inset(AxisZ, 5)
	dim6 := dim.Inset(AxisZ, -5)
	if dim.Left() != 0 {
		t.Errorf("Inset should keep immutable, got left %f", dim.Left())
	}
	if dim.Right() != 20 {
		t.Errorf("Inset should keep immutable, got left %f", dim.Right())
	}
	if dim1.Left() != 5 {
		t.Errorf("Inset failed left should be 5, got %f", dim1.Left())
	}
	if dim1.Right() != 20 {
		t.Errorf("Inset failed left should be 20, got %f", dim1.Right())
	}
	if dim2.Left() != 0 {
		t.Errorf("Inset failed left should be 0, got %f", dim2.Left())
	}
	if dim2.Right() != 15 {
		t.Errorf("Inset failed left should be 15, got %f", dim2.Right())
	}
	if dim3.Top() != 5 {
		t.Errorf("Inset failed left should be 5, got %f", dim3.Top())
	}
	if dim3.Bottom() != 20 {
		t.Errorf("Inset failed left should be 20, got %f", dim3.Bottom())
	}
	if dim4.Top() != 0 {
		t.Errorf("Inset failed left should be 0, got %f", dim4.Top())
	}
	if dim4.Bottom() != 15 {
		t.Errorf("Inset failed left should be 15, got %f", dim4.Bottom())
	}
	if dim5.Front() != 5 {
		t.Errorf("Inset failed left should be 5, got %f", dim5.Front())
	}
	if dim5.Back() != 20 {
		t.Errorf("Inset failed left should be 20, got %f", dim5.Back())
	}
	if dim6.Front() != 0 {
		t.Errorf("Inset failed left should be 0, got %f", dim6.Front())
	}
	if dim6.Back() != 15 {
		t.Errorf("Inset failed left should be 15, got %f", dim6.Back())
	}
}

var aabb3d1 = NewAlignedBox3D(-10, -10, -10, 10, 10, 10)

// 2 instersects with 1
var aabb3d2 = NewAlignedBox3D(-5, 4, -30, -2, 20, 30)

// 3 contains 1
var aabb3d3 = NewAlignedBox3D(-25, -25, -25, 25, 25, 25)

// 4 does not intersect with 1
var aabb3d4 = NewAlignedBox3D(-40, -40, -40, -15, -25, -12)

// Barely intersects with 1
var aabb3d5 = NewAlignedBox3D(10, 10, 10, 20, 20, 20)

// 1 in each direction
var aabb3d6 = NewAlignedBox3D(11, 11, 11, 20, 20, 20)

var aabb2d1 = NewAlignedBox2D(-10, -10, 10, 10)

// 2 instersects with 1
var aabb2d2 = NewAlignedBox2D(-5, 4, -2, 20)

// 3 contains 1
var aabb2d3 = NewAlignedBox2D(-25, -25, 25, 25)

// 4 does not intersect with 1
var aabb2d4 = NewAlignedBox2D(-40, -40, -15, -25)

// Barely intersects with 1
var aabb2d5 = NewAlignedBox2D(10, 10, 20, 20)

// 1 in each direction
var aabb2d6 = NewAlignedBox2D(11, 11, 20, 20)

var distanceTests = []struct {
	d1              *AlignedBox
	d2              *AlignedBox
	distanceSquared float64
}{
	{ // Intersects
		aabb3d1,
		aabb3d2,
		0,
	},
	{ // Contains
		aabb3d1,
		aabb3d3,
		0,
	},
	{
		aabb3d1,
		aabb3d5,
		0,
	},
	{
		aabb3d1,
		aabb3d6,
		3,
	},
	{ //Distance
		aabb3d4,
		aabb3d1,
		254, // 5 * 5 + 15 * 15 + 2 * 2
	},
	{ // Intersects
		aabb2d1,
		aabb2d2,
		0,
	},
	{ // Contains
		aabb2d1,
		aabb2d3,
		0,
	},
	{
		aabb2d1,
		aabb2d5,
		0,
	},
	{
		aabb2d1,
		aabb2d6,
		2,
	},
	{ //Distance
		aabb2d4,
		aabb2d1,
		250, // 5 * 5 + 15 * 15
	},
}

func TestDistance(t *testing.T) {
	for i, args := range distanceTests {
		d := args.d1.DistanceSquared(args.d2)
		d1 := args.d2.DistanceSquared(args.d1)
		if math.Abs(d-args.distanceSquared) > 0.0000001 {
			t.Errorf("Distance %d wrong, expected %f, got %f", i, args.distanceSquared, d)
		}
		if math.Abs(d1-d) > 0.0000001 {
			t.Errorf("Distance %d inverse failed", i)
		}
	}
}
