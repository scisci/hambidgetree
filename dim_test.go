package hambidgetree

import "testing"

var intersectionTests = []struct {
	d1                *Dimension
	d2                *Dimension
	leftIntersection  Extent
	rightIntersection Extent
}{
	{ // First block to the left of second block, same height
		NewDimension2D(0, 0, 100, 100),
		NewDimension2D(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 100),
	},
	{ // First block to the right of second block, same height
		NewDimension2D(100, 0, 200, 100),
		NewDimension2D(0, 0, 100, 100),
		NewExtent(0, 100),
		NewExtent(0, 0),
	},
	{ // First block to the left and shifted up from second block
		NewDimension2D(0, -50, 100, 50),
		NewDimension2D(100, 0, 200, 100),
		NewExtent(-50, -50),
		NewExtent(0, 50),
	},
	{ // First block to the left and shifted down from second block
		NewDimension2D(0, 50, 100, 150),
		NewDimension2D(100, 0, 200, 100),
		NewExtent(50, 50),
		NewExtent(50, 100),
	},
	{ // First block to the left but further than epsilon
		NewDimension2D(0, 0, 50, 100),
		NewDimension2D(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 0),
	},
	{ // First block to the right of second block but further than epsilon, same height
		NewDimension2D(110, 0, 200, 100),
		NewDimension2D(0, 0, 100, 100),
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
	dim := NewDimension3D(0, 0, 0, 20, 20, 20)
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
