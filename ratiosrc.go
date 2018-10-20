package hambidgetree

import (
	exprSolver "github.com/scisci/hambidgetree/expr"
	"math"
	"sort"
	"strconv"
)

// Provides a list of ratios and also the logic for storing those ratios via
// an expression.
type RatioSource interface {
	Ratios() Ratios
	Exprs() Exprs
}

// Creates a ratio source based on a list of expressions.
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

	ratios, err := NewRatios(sortedValues)
	if err != nil {
		return nil, err
	}

	return &basicRatioSource{
		ratios: ratios,
		exprs:  Exprs(sortedExprs),
	}, nil
}

// Creates a ratio source based just on a list of numbers.
func NewBasicRatioSource(values []float64) (RatioSource, error) {
	tmp := make([]float64, len(values))
	copy(tmp, values)
	sort.Float64s(tmp)

	ratios, err := NewRatios(tmp)
	if err != nil {
		return nil, err
	}

	exprs := Exprs(make([]string, len(ratios)))
	for i, value := range ratios {
		exprs[i] = strconv.FormatFloat(value, 'f', -1, 64)
	}

	return &basicRatioSource{
		ratios: ratios,
		exprs:  exprs,
	}, nil
}

// Creates a ratio source based on another ratio source, using some subset.
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
