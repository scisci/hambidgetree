package hambidgetree

import "testing"
import "fmt"

var nearestIndexTests = []struct {
	ratio  float64
	ratios []float64
	index  int
}{
	{
		ratio: 0.5,
		ratios: []float64{
			0.2, 0.8,
		},
		index: 0,
	},
	{
		ratio: 5000.0,
		ratios: []float64{
			0.2, 0.3, 2.0, 1000.0,
		},
		index: 3,
	},
	{
		ratio: 0.65,
		ratios: []float64{
			0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8,
		},
		index: 5,
	},
	{
		ratio: 0.25,
		ratios: []float64{
			0.1, 0.2, 0.9, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8,
		},
		index: 1,
	},
	{
		ratio: 1.5,
		ratios: []float64{
			5.0, 6.0, 18.45,
		},
		index: 0,
	},
}

func TestFindNearestIndex(t *testing.T) {
	for i, test := range nearestIndexTests {
		ratioSource, err := NewBasicRatioSource(test.ratios)
		if err != nil {
			t.Errorf("Error creating ratio source %v", err)
		}
		ratios := ratioSource.Ratios()
		index := FindClosestIndex(ratios, test.ratio, 0.0000001)

		if index != test.index {
			t.Errorf("nearest index test %d failed, expected %d, got %d", i, test.index, index)
		}

		inverseIndex := FindClosestInverseIndex(ratios, 1/test.ratio, 0.0000001)
		if inverseIndex != test.index {
			t.Errorf("nearest inverse index test %d failed, expected %d, got %d", i, test.index, index)
		}
	}
}

func TestFindInverseRatioIndex(t *testing.T) {
	ratioSource, err := NewBasicRatioSource([]float64{0.5, 1.0, 1.5, 2.0, 8.3})
	if err != nil {
		t.Errorf("Error creating ratio source %v", err)
	}
	ratios := ratioSource.Ratios()
	index := FindInverseRatioIndex(ratios, 0, 0.0000001)
	if index != 3 {
		t.Errorf("inverse ratio index should be 3, got %d", index)
	}

	index = FindInverseRatioIndex(ratios, 3, 0.0000001)
	if index != 0 {
		t.Errorf("inverse ratio index should be 0, got %d", index)
	}
}

func TestSubset(t *testing.T) {
	ratioSource, err := NewBasicRatioSource([]float64{5.0, 288.04, 7.43, 2828.18, 3.482})
	if err != nil {
		t.Errorf("Error creating ratio source %v", err)
	}
	ratioSourceSubset, err := NewRatioSourceSubset(ratioSource, []float64{5.0, 7.43, 3.482}, 0.0000001)
	if err != nil {
		t.Errorf("failed to create ratio subset %v", err)
	}

	ratios := ratioSourceSubset.Ratios()

	if l := len(ratios); l != 3 {
		t.Errorf("subset length should be 3, got %d", l)
	}

	if v := ratios[0]; v != 3.482 {
		t.Errorf("first subset item should be 3.482, got %f", v)
	}

	if v := ratios[1]; v != 5.0 {
		t.Errorf("second subset item should be 5.0, got %f", v)
	}

	if v := ratios[2]; v != 7.43 {
		t.Errorf("third subset item should be 7.43, got %f", v)
	}
}

func TestComplements(t *testing.T) {
	ratioSource, err := NewBasicRatioSource([]float64{0.25, 0.5, 0.75, 1.0, 1.3333333, 2.0, 4.0})
	if err != nil {
		t.Errorf("Error creating ratio source %v", err)
	}
	complements, err := NewComplements(ratioSource.Ratios(), 0.0000001)
	if err != nil {
		t.Errorf("complements should not have error")
	}

	if len(complements) != len(ratioSource.Ratios()) {
		t.Errorf("complements length should match ratios length")
	}

	fmt.Println(complements)
}
