package hambidgetree

import "math"
import "sort"
import "strconv"

type Complements [][]Split
type Ratios interface {
	Len() int
	At(index int) float64
}

type TreeRatios interface {
	Complements() Complements // Returns complements
	Ratios() Ratios           // Returns ratios sorted ascending
}

type rawRatios []float64

func (ratios rawRatios) Len() int {
	return len(ratios)
}

func (ratios rawRatios) At(index int) float64 {
	return ratios[index]
}

type ratioSubset struct {
	ratios  Ratios
	indexes []int
}

func (ratios *ratioSubset) Len() int {
	return len(ratios.indexes)
}

func (ratios *ratioSubset) At(index int) float64 {
	return ratios.ratios.At(ratios.indexes[index])
}

type treeRatios struct {
	ratios      Ratios
	complements Complements
}

func NewTreeRatios(ratios Ratios, epsilon float64) TreeRatios {
	return &treeRatios{
		ratios:      ratios,
		complements: NewComplements(ratios, epsilon),
	}
}

func (ratios *treeRatios) Ratios() Ratios {
	return ratios.ratios
}

func (ratios *treeRatios) Complements() Complements {
	return ratios.complements
}

func NewRatios(values []float64) Ratios {
	ratios := make([]float64, len(values))
	copy(ratios, values)
	sort.Float64s(ratios)
	return rawRatios(ratios)
}

func NewRatiosSubset(ratios Ratios, values []float64, epsilon float64) Ratios {
	subset := &ratioSubset{
		ratios:  ratios,
		indexes: nil,
	}

	for _, value := range values {
		for j := 0; j < ratios.Len(); j++ {
			if math.Abs(ratios.At(j)-value) < epsilon {
				subset.indexes = append(subset.indexes, j)
				break
			}
		}
	}

	sort.Ints(subset.indexes)
	return subset
}

func NewComplements(ratios Ratios, epsilon float64) Complements {
	n := ratios.Len()
	complements := make([][]Split, n)

	for i := 0; i < n; i++ {
		ratio := ratios.At(i)

		// Try to split the width, in the ratio array the height is always considered
		// to be unity
		for j := 0; j < n; j++ {
			if ratios.At(j) < ratio-epsilon {
				for k := j; k < n; k++ {
					if math.Abs(ratio-ratios.At(j)-ratios.At(k)) < epsilon {
						left, right := j, k
						if left > right {
							left, right = right, left
						}
						complements[i] = append(complements[i], NewVerticalSplit(left, right))
						break
					}
				}
			}
		}

		// Now try to split the height, we need to invert the ratio, since all
		// ratios are setup based on unity height,
		// So if we have a ratio that is 0.5, and we want to split its height, we
		// consider the ratio as 1.0 / 0.5, or 2.0, then we figure out how can you
		// split a ratio of 2.0, well with 2 1.0s, etc.
		ratio = 1.0 / ratio
		for j := 0; j < n; j++ {
			if ratios.At(j) < ratio-epsilon {
				for k := j; k < n; k++ {
					if math.Abs(ratio-ratios.At(j)-ratios.At(k)) < epsilon {
						// The ratios here won't actually be j and k, they will be the inverses
						// Because they are vertically stacked and height is always
						// considered 1
						invJ := FindInverseRatioIndex(ratios, j, epsilon)
						invK := FindInverseRatioIndex(ratios, k, epsilon)

						if invJ < 0 || invK < 0 {
							panic("inverse ratio lookup failed " + strconv.Itoa(j) + ":" + strconv.Itoa(invJ) + ", " + strconv.Itoa(k) + ":" + strconv.Itoa(invK))
						}

						top, bot := invJ, invK
						if top > bot {
							top, bot = bot, top
						}
						complements[i] = append(complements[i], NewHorizontalSplit(top, bot))
						break
					}
				}
			}
		}
	}

	return Complements(complements)
}

func (c Complements) String() string {
	str := ""
	for i := range c {
		str = str + "\n" + strconv.Itoa(i) + ":\n"
		for j := range c[i] {
			str = str + " " + c[i][j].String()
		}
	}

	return str
}

// Binary search to find closest index, must provided an epsilon for float precision
// errors, if there are two that are the same distance, the smaller index wins.
func FindClosestIndex(ratios Ratios, ratio, epsilon float64) int {
	closestDist := math.MaxFloat64
	closestIndex := -1

	//loops := 0

	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, ratios.Len()
	for i < j {
		//loops = loops + 1
		h := i + (j-i)/2 // avoid overflow when computing h

		dist := ratio - ratios.At(h)

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

func FindClosestInverseIndex(ratios Ratios, ratio, epsilon float64) int {
	return FindClosestIndex(ratios, 1.0/ratio, epsilon)
}

func FindInverseRatioIndex(ratios Ratios, index int, epsilon float64) int {
	inverseRatio := 1.0 / ratios.At(index)
	closestIndex := FindClosestIndex(ratios, inverseRatio, epsilon)
	if closestIndex >= 0 {
		if math.Abs(ratios.At(closestIndex)-inverseRatio) < epsilon {
			return closestIndex
		}
	}

	return -1
}
