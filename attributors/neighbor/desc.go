package neighbor

import (
	"github.com/scisci/hambidgetree/attributors"
)

func (attributor *HasNeighborAttributor) Name() string {
	return "HasNeighbor"
}

func (attributor *HasNeighborAttributor) Description() string {
	return "This attributor finds all of the leaves that have at least " +
		"one neighbor touching on either side, then chooses a random one and " +
		"marks it. Once marked, the leaf is considered 'deleted.'"
}

func (attributor *HasNeighborAttributor) Parameters(f attributors.ParameterFormatType) map[string]interface{} {
	return map[string]interface{}{
		"Max Marks": attributor.MaxMarks,
		"Dimension": attributor.Dimension,
		"Seed":      attributor.Seed,
	}
}
