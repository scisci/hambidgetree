package stepped

import "math"

func ValueToSteps(value float64, max int) int {
	stepped := int(math.Floor(value*float64(max) + 0.5))
	if stepped < 0 {
		stepped = 0
	} else if stepped > max {
		stepped = max
	}
	return stepped
}

func StepsToValue(steps int, max int) float64 {
	return float64(steps) / float64(max)
}
