package hambidgetree

import (
	exprSolver "github.com/scisci/hambidgetree/expr"
	"math"
	"sort"
	"strconv"
)

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
