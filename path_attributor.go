package hambidgetree

import (
	"bytes"
	"fmt"
	gpath "gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"math"
	"math/rand"
)

// Create a path finder filter
// Chooses one or more paths through a tree and marks those nodes

// Paths could be from node to node, edge to edge, etc.
//
type EdgeName int

const (
	EdgeNameLeft   EdgeName = 0
	EdgeNameRight  EdgeName = iota
	EdgeNameTop    EdgeName = iota
	EdgeNameBottom EdgeName = iota
	EdgeNameFront  EdgeName = iota
	EdgeNameBack   EdgeName = iota
)

func (name EdgeName) String() string {
	switch name {
	case EdgeNameLeft:
		return "Left"
	case EdgeNameRight:
		return "Right"
	case EdgeNameTop:
		return "Top"
	case EdgeNameBottom:
		return "Bottom"
	case EdgeNameFront:
		return "Front"
	case EdgeNameBack:
		return "Back"
	}

	return "Unknown"
}

func (name EdgeName) Index() int64 {
	return int64(name)
}

// TODO: These won't be right for trees created with different aspect ratios!
func (name EdgeName) Dimension3D() *Dimension {
	switch name {
	case EdgeNameLeft:
		return NewDimension3D(0, 0, 0, 0, 1, 1)
	case EdgeNameRight:
		return NewDimension3D(1, 0, 0, 1, 1, 1)
	case EdgeNameTop:
		return NewDimension3D(0, 0, 0, 1, 0, 1)
	case EdgeNameBottom:
		return NewDimension3D(0, 1, 0, 1, 1, 1)
	case EdgeNameFront:
		return NewDimension3D(0, 0, 0, 1, 1, 0)
	case EdgeNameBack:
		return NewDimension3D(0, 0, 1, 1, 1, 1)
	}

	panic("unknown edge name!")
	return nil
}

type EdgePath struct {
	From EdgeName
	To   EdgeName
}

type edgeNode struct {
	name EdgeName
	id   int64
	dim  *Dimension
}

func createEdgeNodes(idOffset int64) []*edgeNode {
	return []*edgeNode{
		&edgeNode{EdgeNameLeft, idOffset + EdgeNameLeft.Index(), EdgeNameLeft.Dimension3D()},
		&edgeNode{EdgeNameRight, idOffset + EdgeNameRight.Index(), EdgeNameRight.Dimension3D()},
		&edgeNode{EdgeNameTop, idOffset + EdgeNameTop.Index(), EdgeNameTop.Dimension3D()},
		&edgeNode{EdgeNameBottom, idOffset + EdgeNameBottom.Index(), EdgeNameBottom.Dimension3D()},
		&edgeNode{EdgeNameFront, idOffset + EdgeNameFront.Index(), EdgeNameFront.Dimension3D()},
		&edgeNode{EdgeNameBack, idOffset + EdgeNameBack.Index(), EdgeNameBack.Dimension3D()},
	}
}

type edgePathParams []EdgePath

func (paths edgePathParams) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[")
	for i, p := range paths {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(p.From.String())
		buf.WriteString("->")
		buf.WriteString(p.To.String())
	}
	buf.WriteString("]")
	return buf.String()
}

var OnPathAttr = "onPath"
var OnPathValue = "true"

type EdgePathAttributor struct {
	Paths    []EdgePath
	Chaos    int
	MaxChaos int
	Seed     int64
}

const edgePathMaxChaos = 100

func NewEdgePathAttributor(paths []EdgePath, seed int64, chaos float64) *EdgePathAttributor {
	return &EdgePathAttributor{
		Paths:    paths,
		Seed:     seed,
		Chaos:    valueToSteps(chaos, edgePathMaxChaos),
		MaxChaos: edgePathMaxChaos,
	}
}

func (attributor *EdgePathAttributor) Name() string {
	return "EdgePath"
}

func (attributor *EdgePathAttributor) Description() string {
	return "This attributor finds a path between the edges and marks each item along that path with the attribute."
}

func (attributor *EdgePathAttributor) Parameters(f ParameterFormatType) map[string]interface{} {
	return map[string]interface{}{
		"Paths":    edgePathParams(attributor.Paths),
		"Seed":     attributor.Seed,
		"Chaos":    attributor.Chaos,
		"MaxChaos": attributor.MaxChaos,
	}
}

func (attributor *EdgePathAttributor) AddAttributes(tree *Tree, attrs *NodeAttributer) error {
	rand.Seed(attributor.Seed)
	//epsilon := 0.0000001

	dimensionLookup := NewNodeDimensionMap(tree, NewVector(0, 0, 0), 1.0)
	matrix, err := BuildAdjacencyMatrix(tree, dimensionLookup)
	if err != nil {
		return err
	}
	numLeaves := len(matrix)
	fmt.Printf("working with %d leaves\n", numLeaves)

	maxID := int64(0)
	for leafID, _ := range matrix {
		if int64(leafID) > maxID {
			maxID = int64(leafID)
		}
	}

	edges := createEdgeNodes(maxID + 1)

	graph := simple.NewWeightedUndirectedGraph(0, math.Inf(1))

	// Value 0 to 1, 0 is a shortest path (distance 1) 1, is complete noise (distance is rand * numleaves)
	chaos := stepsToValue(attributor.Chaos, attributor.MaxChaos)

	for leafID, neighbors := range matrix {
		fmt.Printf("leaf %d has %d neighbors\n", leafID, len(neighbors))

		for _, neighbor := range neighbors {
			if graph.HasEdgeBetween(simple.Node(leafID), simple.Node(neighbor.ID())) {
				continue
			}

			randWeight := 1 + chaos*rand.Float64()*float64(numLeaves)
			graph.SetWeightedEdge(&simple.WeightedEdge{F: simple.Node(leafID), T: simple.Node(neighbor.ID()), W: randWeight})
			fmt.Printf("%d to %d = %f\n", leafID, neighbor.ID(), randWeight)

		}

		// Also need to check this neighbor against the edges
		dim, err := dimensionLookup.Dimension(leafID)
		if err != nil {
			return err
		}

		for _, e := range edges {
			if e.dim.DistanceSquared(dim) < 0.0000001 {
				edgeWeight := 1.0
				graph.SetWeightedEdge(&simple.WeightedEdge{F: simple.Node(leafID), T: simple.Node(e.id), W: edgeWeight})
			}
		}
	}

	for _, path := range attributor.Paths {
		var fromNode *edgeNode
		var toNode *edgeNode
		for _, edge := range edges {
			if edge.name == path.From {
				fromNode = edge
			}
			if edge.name == path.To {
				toNode = edge
			}
		}
		if fromNode == nil || toNode == nil {
			panic("path ends not found")
		}

		shortest := gpath.DijkstraFrom(simple.Node(fromNode.id), graph)
		p, w := shortest.To(simple.Node(toNode.id))
		fmt.Printf("path weight %f\n", w)

		for i, graphNode := range p {
			if i > 0 && i < len(p)-1 {
				attrs.SetAttribute(NodeID(graphNode.ID()), OnPathAttr, OnPathValue)
			}
		}
	}

	return nil
}
