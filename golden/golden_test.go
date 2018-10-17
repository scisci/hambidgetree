package golden_test

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/golden"
	"testing"
)

func TestGolden(t *testing.T) {
	r1 := htree.NewBasicRatioSource(golden.Floats)
	r2, err := htree.NewExprRatioSource(golden.Exprs)
	if err != nil {
		t.Errorf("Golden ratio expr error %v", err)
	}

	r1f := r1.Ratios()
	r2f := r2.Ratios()

	if len(r1f) != len(r2f) {
		t.Errorf("Golden ratios lengths don't match %d & %d", len(r1f), len(r2f))
	}

	for i := 0; i < len(r1f); i++ {
		v1 := r1f[i]
		v2 := r2f[i]

		if v1 != v2 {
			t.Errorf("Golden value at %d, don't match %f and %f", i, v1, v2)
		}
	}
}
