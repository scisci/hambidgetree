package generators

import (
	htree "github.com/scisci/hambidgetree"
)

type ParameterFormatType int

const (
	ParameterFormatTypeVerbose ParameterFormatType = 0
	ParameterFormatTypeConcise ParameterFormatType = 1
)

type TreeGenerator interface {
	Name() string
	Description() string
	Parameters(f ParameterFormatType) map[string]interface{}
	Generate() (htree.Tree, error)
}
