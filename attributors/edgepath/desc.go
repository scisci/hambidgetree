package edgepath

import (
	"github.com/scisci/hambidgetree/attributors"
)

func (attributor *EdgePathAttributor) Name() string {
	return "EdgePath"
}

func (attributor *EdgePathAttributor) Description() string {
	return "This attributor finds a path between the edges and marks each item along that path with the attribute."
}

func (attributor *EdgePathAttributor) Parameters(f attributors.ParameterFormatType) map[string]interface{} {
	return map[string]interface{}{
		"Paths":    edgePathParams(attributor.Paths),
		"Seed":     attributor.Seed,
		"Chaos":    attributor.Chaos,
		"MaxChaos": attributor.MaxChaos,
	}
}
