package hambidgetree

type NodeAttributes interface {
	Attribute(node *Node, key string) (string, error)
}

type NodeAttributer struct {
	attrs map[NodeID]map[string]string
}

func NewNodeAttributer() *NodeAttributer {
	return &NodeAttributer{
		attrs: make(map[NodeID]map[string]string),
	}
}

func (attributer *NodeAttributer) SetAttribute(node *Node, key, value string) {
	attrs, ok := attributer.attrs[node.id]
	if !ok {
		attrs = make(map[string]string)
		attributer.attrs[node.id] = attrs
	}

	attrs[key] = value
}

func (attributer *NodeAttributer) Attribute(node *Node, key string) (string, error) {
	attrs, ok := attributer.attrs[node.id]
	if !ok {
		return "", ErrNotFound
	}

	value, ok := attrs[key]
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}
