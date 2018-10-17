package hambidgetree

import (
	"errors"
	"math"
	"strconv"
)

var ErrMissingInverse = errors.New("Missing inverse")

func NewComplements(ratios Ratios, epsilon float64) (Complements, error) {
	for i := 0; i < len(ratios); i++ {
		if FindInverseRatioIndex(ratios, i, epsilon) == -1 {
			return nil, ErrMissingInverse
		}
	}

	n := len(ratios)
	complements := make([][]Split, n)

	for i := 0; i < n; i++ {
		ratio := ratios[i]

		// Try to split the width, in the ratio array the height is always considered
		// to be unity
		for j := 0; j < n; j++ {
			if ratios[j] < ratio-epsilon {
				for k := j; k < n; k++ {
					if math.Abs(ratio-ratios[j]-ratios[k]) < epsilon {
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
			if ratios[j] < ratio-epsilon {
				for k := j; k < n; k++ {
					if math.Abs(ratio-ratios[j]-ratios[k]) < epsilon {
						// The ratios here won't actually be j and k, they will be the inverses
						// Because they are vertically stacked and height is always
						// considered 1
						invJ := FindInverseRatioIndex(ratios, j, epsilon)
						invK := FindInverseRatioIndex(ratios, k, epsilon)

						if invJ < 0 || invK < 0 {
							panic("inverse ratio lookup failed " + strconv.Itoa(j) + ":" +
								strconv.Itoa(invJ) + ", " + strconv.Itoa(k) + ":" + strconv.Itoa(invK))
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

	return Complements(complements), nil
}
