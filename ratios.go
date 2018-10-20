package hambidgetree

import (
	"errors"
	"math"
)

var ErrRatiosUnordered = errors.New("Ratios unordered")
var ErrRatiosContainsDuplicates = errors.New("Ratios contain duplicates")
var ErrRatioNotFound = errors.New("Ratio not found")

const RatioIndexUndefined = -1

// Ratios are guaranteed to be sorted ascending and contain no duplicates
type Ratios []float64

// Exprs guaranteed to be sorted ascending and contain no duplicates
type Exprs []string

// A list of possible splits mapped by index to a ratios array
type Complements [][]Split

func IsRatioIndexDefined(index int) bool {
	return index > RatioIndexUndefined
}

// Creates ratios from list of values, enforces that values are sorted
// ascending and no duplicates, otherwise an error is thrown.
func NewRatios(values []float64) (Ratios, error) {
	n := len(values)

	for i := 1; i < n; i++ {
		if values[i-1] > values[i] {
			return nil, ErrRatiosUnordered
		}

		if values[i-1] == values[i] {
			return nil, ErrRatiosContainsDuplicates
		}
	}

	return Ratios(values), nil
}

// Returns the parameterized height of a ratio if it is contained within another
// ratio. i.e. A ratio of 2:1 within a ratio of 1:2 has a normal height of 0.25
func RatioNormalHeight(containerRatio, ratio float64) float64 {
	return containerRatio / ratio
}

// Returns the parameterized width of a ratio if it is contained within another
// ratio. i.e. A ratio of 1:2 within a ratio of 2:1 has a normal width of 0.25
func RatioNormalWidth(containerRatio, ratio float64) float64 {
	return ratio / containerRatio
}

// Given a sorted list of values, the epsilon value is some function of the
// minimum distance between two of the values. Technically we should only have
// to divide this by 2, but we do it by 1000 just for fun.
func CalculateRatiosEpsilon(ratios Ratios) float64 {
	minDist := math.MaxFloat64
	n := len(ratios)

	if n > 1 {
		lastVal := ratios[0]
		for i := 1; i < n; i++ {
			val := ratios[i]
			if val-lastVal < minDist {
				minDist = val - lastVal
			}
			lastVal = val
		}
	}

	return minDist / 1000.0
}

func FindIndexesWithMissingInverses(ratios Ratios, epsilon float64) []int {
	var indexes []int

	for i := 0; i < len(ratios); i++ {
		if FindInverseRatioIndex(ratios, i, epsilon) == -1 {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

// Binary search to find closest index, must provided an epsilon for float precision
// errors, if there are two that are the same distance, the smaller index wins.
func FindClosestIndex(ratios Ratios, ratio, epsilon float64) int {
	closestDist := math.MaxFloat64
	closestIndex := RatioIndexUndefined

	//loops := 0

	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, len(ratios)
	for i < j {
		//loops = loops + 1
		h := i + (j-i)/2 // avoid overflow when computing h

		dist := ratio - ratios[h]

		// i ≤ h < j
		if dist > 0 {
			i = h + 1 // preserves f(i-1) == false
		} else if dist < 0 {
			dist = -dist
			j = h // preserves f(j) == true
		} else {
			closestIndex = h
			break
		}

		if dist < closestDist-epsilon || (dist < closestDist+epsilon && h < closestIndex) {
			closestDist = dist
			closestIndex = h
		}
	}

	return closestIndex
}

func FindClosestIndexWithinRange(ratios Ratios, ratio, epsilon float64) int {
	index := FindClosestIndex(ratios, ratio, epsilon)
	if index < 0 {
		return RatioIndexUndefined
	}

	dist := ratio - ratios[index]
	if dist < -epsilon || dist > epsilon {
		return RatioIndexUndefined
	}

	return index
}

func FindClosestInverseIndex(ratios Ratios, ratio, epsilon float64) int {
	return FindClosestIndex(ratios, 1.0/ratio, epsilon)
}

func FindInverseRatioIndex(ratios Ratios, index int, epsilon float64) int {
	inverseRatio := 1.0 / ratios[index]
	closestIndex := FindClosestIndex(ratios, inverseRatio, epsilon)
	if closestIndex >= 0 {
		if math.Abs(ratios[closestIndex]-inverseRatio) < epsilon {
			return closestIndex
		}
	}

	return RatioIndexUndefined
}
