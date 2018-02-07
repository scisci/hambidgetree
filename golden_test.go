package hambidgetree

import "testing"

func TestGolden(t *testing.T) {
	r1 := NewGoldenRatios()
	r2 := NewGoldenRatiosRaw()

	if r1.Len() != r2.Len() {
		t.Errorf("Golden ratios lengths don't match %d & %d", r1.Len(), r2.Len())
	}

	for i := 0; i < r1.Len(); i++ {
		v1 := r1.At(i)
		v2 := r2.At(i)

		if v1 != v2 {
			t.Errorf("Golden value at %d, don't match %f and %f", i, v1, v2)
		}
	}
}
