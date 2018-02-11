package factory

import (
	"encoding/json"
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/simple"
)

type jsonWrapper struct {
	Type    string          `json:"type"`
	Version int             `json:"version"`
	Tree    json.RawMessage `json:"tree"`
}

func UnmarshalJSON(data []byte) (htree.ImmutableTree, error) {
	var wrapper *jsonWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	if wrapper.Type == "simple" {
		return simple.UnmarshalJSON(wrapper.Version, wrapper.Tree)
	}

	return nil, fmt.Errorf("Unknown type %s", wrapper.Type)
}

func MarshalJSON(tree htree.ImmutableTree) ([]byte, error) {
	_, ok := tree.(*simple.Tree)
	if !ok {
		return nil, fmt.Errorf("Unknown tree type!")
	}

	simpleData, err := json.Marshal(tree)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jsonWrapper{
		Type:    "simple",
		Version: simple.JSONVersion,
		Tree:    simpleData,
	})
}
