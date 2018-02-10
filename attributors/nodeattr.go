package attributors

import (
	htree "github.com/scisci/hambidgetree"
)

type NodeAttributes interface {
	Attribute(id htree.NodeID, key string) (string, error)
}

type NodeAttributer struct {
	attrs map[htree.NodeID]map[string]string
}

func NewNodeAttributer() *NodeAttributer {
	return &NodeAttributer{
		attrs: make(map[htree.NodeID]map[string]string),
	}
}

func (attributer *NodeAttributer) SetAttribute(id htree.NodeID, key, value string) {
	attrs, ok := attributer.attrs[id]
	if !ok {
		attrs = make(map[string]string)
		attributer.attrs[id] = attrs
	}

	attrs[key] = value
}

func (attributer *NodeAttributer) Attribute(id htree.NodeID, key string) (string, error) {
	attrs, ok := attributer.attrs[id]
	if !ok {
		return "", ErrNotFound
	}

	value, ok := attrs[key]
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}
