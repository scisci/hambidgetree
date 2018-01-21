package hambidgetree

import "testing"

func TestMapSteps(t *testing.T) {
	s := stepsToValue(512, 1024)
	if s != 0.5 {
		t.Errorf("Value should be 0.5, got %f", s)
	}

	for i := 0; i <= 1024; i++ {
		s = stepsToValue(i, 1024)
		v := valueToSteps(s, 1024)
		if v != i {
			t.Errorf("Steps should equal %d, got %d", i, v)
		}
	}
}
