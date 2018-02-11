package attributors

import (
	"errors"
	htree "github.com/scisci/hambidgetree"
)

var ErrNotFound = errors.New("Not Found")

type ParameterFormatType int

const (
	ParameterFormatTypeVerbose ParameterFormatType = 0
	ParameterFormatTypeConcise ParameterFormatType = 1
)

type TreeAttributor interface {
	Name() string
	Description() string
	Parameters(f ParameterFormatType) map[string]interface{}
	AddAttributes(tree htree.ImmutableTree, attrs *NodeAttributer) error
}
