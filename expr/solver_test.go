package expr

import "testing"
import "math"

var tests = []struct {
	Expr   string
	Result float64
}{
	{"1/2", 0.5},
	{"0.125", 0.125},
	{"1 / (SQRT(5) + 5)", 1 / (math.Sqrt(5) + 5)},
	{"2/(SQRT(5)*3+7)", 2 / (math.Sqrt(5)*3 + 7)},
	{"SQRT(5) / (4+4*SQRT(5))", math.Sqrt(5) / (4 + 4*math.Sqrt(5))},
	{"4/(9+SQRT(5))", 0.3559964222368531738732013858562486},
	{"SQRT(5)/2 + 0.5", math.Phi},
	{"PHI", (math.Sqrt(5) + 1) / 2},
}

func TestSolver(t *testing.T) {
	for i, test := range tests {
		res, err := Solve(test.Expr)
		if err != nil {
			t.Errorf("Failed to parse %v", err)
		}
		if res != test.Result {
			t.Errorf("Test %d failed, expected %.20f, got %.20f", i, test.Result, res)
		}
	}
}
