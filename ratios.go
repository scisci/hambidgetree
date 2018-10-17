package hambidgetree

import (
	"errors"
	exprSolver "github.com/scisci/hambidgetree/expr"
	"math"
	"sort"
	"strconv"
)

var ErrRatiosUnordered = errors.New("Ratios unordered")
var ErrRatiosContainsDuplicates = errors.New("Ratios contain duplicates")
var ErrRatioNotFound = errors.New("Ratio not found")

const RatioIndexUndefined = -1

type Ratios []float64
type Exprs []string
type Complements [][]Split

type RatioSource interface {
	Ratios() Ratios
	Exprs() Exprs
}

func NewExprRatioSource(exprs []string) (RatioSource, error) {
	var tmp exprValues
	for _, expr := range exprs {
		value, err := exprSolver.Solve(expr)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, exprValue{expr: expr, value: value})
	}
	sort.Sort(tmp)

	sortedValues := make([]float64, len(tmp))
	sortedExprs := make([]string, len(tmp))
	for i, exprValue := range tmp {
		sortedValues[i] = exprValue.value
		sortedExprs[i] = exprValue.expr
	}

	return &basicRatioSource{
		ratios: Ratios(sortedValues),
		exprs:  Exprs(sortedExprs),
	}, nil
}

func NewBasicRatioSource(values []float64) RatioSource {
	tmp := make([]float64, len(values))
	copy(tmp, values)
	sort.Float64s(tmp)

	ratios, err := NewRatios(tmp)
	if err != nil {
		panic(err)
	}

	exprs := Exprs(make([]string, len(ratios)))
	for i, value := range ratios {
		exprs[i] = strconv.FormatFloat(value, 'f', -1, 64)
	}

	return &basicRatioSource{
		ratios: ratios,
		exprs:  exprs,
	}
}

type basicRatioSource struct {
	ratios Ratios
	exprs  []string
}

func (basicRatioSource *basicRatioSource) Ratios() Ratios {
	return basicRatioSource.ratios
}

func (basicRatioSource *basicRatioSource) Exprs() Exprs {
	return basicRatioSource.exprs
}

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

// A pairing of an expression and value useful for creating a ratio source from
// expressions.
type exprValue struct {
	value float64
	expr  string
}

type exprValues []exprValue

func (a exprValues) Len() int           { return len(a) }
func (a exprValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a exprValues) Less(i, j int) bool { return a[i].value < a[j].value }

// Concrete struct for implementing the RatioSourceSubSet interface
type ratioSourceSubset struct {
	ratioSource RatioSource
	ratios      Ratios
	exprs       Exprs
	indexes     []int
}

func (ratioSourceSubset *ratioSourceSubset) Ratios() Ratios {
	return ratioSourceSubset.ratios
}

func (ratioSourceSubset *ratioSourceSubset) Exprs() Exprs {
	return ratioSourceSubset.exprs
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

func NewRatioSourceSubset(ratioSource RatioSource, values []float64, epsilon float64) (RatioSource, error) {
	tmp := make([]float64, len(values))
	copy(tmp, values)
	sort.Float64s(tmp)

	allRatios := ratioSource.Ratios()
	allExprs := ratioSource.Exprs()

	subExprs := Exprs(make([]string, len(tmp)))
	subRatios, err := NewRatios(tmp)
	if err != nil {
		panic(err)
	}

	var indexes []int
	for i, value := range subRatios {
		found := false
		for j := 0; j < len(allRatios); j++ {
			if math.Abs(allRatios[j]-value) < epsilon {
				found = true
				indexes = append(indexes, j)
				subExprs[i] = allExprs[j]
				break
			}
		}
		if !found {
			return nil, ErrRatioNotFound
		}
	}

	return &ratioSourceSubset{
		ratioSource: ratioSource,
		ratios:      subRatios,
		exprs:       subExprs,
		indexes:     indexes,
	}, nil
}

// Binary search to find closest index, must provided an epsilon for float precision
// errors, if there are two that are the same distance, the smaller index wins.
func FindClosestIndex(ratios Ratios, ratio, epsilon float64) int {
	closestDist := math.MaxFloat64
	closestIndex := -1

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
		return -1
	}

	dist := ratio - ratios[index]
	if dist < -epsilon || dist > epsilon {
		return -1
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

	return -1
}
