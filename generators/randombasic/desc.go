package randombasic

import (
	"github.com/scisci/hambidgetree/print"
	"github.com/scisci/hambidgetree/generators"
	"strconv"
)

func (gen *RandomBasicTreeGenerator) Name() string {
	return "Random Basic"
}

func (gen *RandomBasicTreeGenerator) Description() string {
	return "This algorithm can be set to generate a given number of leaves. It begins with a single leaf of a given ratio. Until it reaches the desired number of leaves, it selects a leaf at random and splits it."
}

func (gen *RandomBasicTreeGenerator) Parameters(f generators.ParameterFormatType) map[string]interface{} {
	if f == generators.ParameterFormatTypeConcise {
		return map[string]interface{}{
			"# Leaves": gen.NumLeaves,
			"Seed":     gen.Seed,
		}
	}

	return map[string]interface{}{
		"Ratios":               print.PrintRatios(gen.RatioSource),
		"Container Ratio (XY)": strconv.FormatFloat(gen.XYRatio, 'f', 4, 64),
		"Container Ratio (ZY)": strconv.FormatFloat(gen.ZYRatio, 'f', 4, 64),
		"Number of Leaves":     gen.NumLeaves,
		"Random Seed":          gen.Seed,
	}
}
