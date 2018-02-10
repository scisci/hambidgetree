package stepped

import "testing"

func TestMapSteps(t *testing.T) {
	s := StepsToValue(512, 1024)
	if s != 0.5 {
		t.Errorf("Value should be 0.5, got %f", s)
	}

	for i := 0; i <= 1024; i++ {
		s = StepsToValue(i, 1024)
		v := ValueToSteps(s, 1024)
		if v != i {
			t.Errorf("Steps should equal %d, got %d", i, v)
		}
	}
}
