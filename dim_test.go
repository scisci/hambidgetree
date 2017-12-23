package hambidgetree

import "testing"

var intersectionTests = []struct {
	d1                *Dimension
	d2                *Dimension
	leftIntersection  Extent
	rightIntersection Extent
}{
	{ // First block to the left of second block, same height
		NewDimension(0, 0, 100, 100),
		NewDimension(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 100),
	},
	{ // First block to the right of second block, same height
		NewDimension(100, 0, 200, 100),
		NewDimension(0, 0, 100, 100),
		NewExtent(0, 100),
		NewExtent(0, 0),
	},
	{ // First block to the left and shifted up from second block
		NewDimension(0, -50, 100, 50),
		NewDimension(100, 0, 200, 100),
		NewExtent(-50, -50),
		NewExtent(0, 50),
	},
	{ // First block to the left and shifted down from second block
		NewDimension(0, 50, 100, 150),
		NewDimension(100, 0, 200, 100),
		NewExtent(50, 50),
		NewExtent(50, 100),
	},
	{ // First block to the left but further than epsilon
		NewDimension(0, 0, 50, 100),
		NewDimension(100, 0, 200, 100),
		NewExtent(0, 0),
		NewExtent(0, 0),
	},
	{ // First block to the right of second block but further than epsilon, same height
		NewDimension(110, 0, 200, 100),
		NewDimension(0, 0, 100, 100),
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
