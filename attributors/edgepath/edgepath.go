package edgepath

import (
	"bytes"
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/algo"
	"github.com/scisci/hambidgetree/attributors"
	"github.com/scisci/hambidgetree/stepped"
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
func (name EdgeName) Dimension3D() *htree.Dimension {
	switch name {
	case EdgeNameLeft:
		return htree.NewDimension3D(0, 0, 0, 0, 1, 1)
	case EdgeNameRight:
		return htree.NewDimension3D(1, 0, 0, 1, 1, 1)
	case EdgeNameTop:
		return htree.NewDimension3D(0, 0, 0, 1, 0, 1)
	case EdgeNameBottom:
		return htree.NewDimension3D(0, 1, 0, 1, 1, 1)
	case EdgeNameFront:
		return htree.NewDimension3D(0, 0, 0, 1, 1, 0)
	case EdgeNameBack:
		return htree.NewDimension3D(0, 0, 1, 1, 1, 1)
	}

	panic("unknown edge name!")
	return nil
}

type EdgePath struct {
	From EdgeName
	To   EdgeName
}

type edgeNode struct {
	name      EdgeName
	id        int64
	dim       *htree.Dimension
	neighbors []htree.NodeID
}

func createEdgeNodes(idOffset int64) []*edgeNode {
	return []*edgeNode{
		&edgeNode{EdgeNameLeft, idOffset + EdgeNameLeft.Index(), EdgeNameLeft.Dimension3D(), nil},
		&edgeNode{EdgeNameRight, idOffset + EdgeNameRight.Index(), EdgeNameRight.Dimension3D(), nil},
		&edgeNode{EdgeNameTop, idOffset + EdgeNameTop.Index(), EdgeNameTop.Dimension3D(), nil},
		&edgeNode{EdgeNameBottom, idOffset + EdgeNameBottom.Index(), EdgeNameBottom.Dimension3D(), nil},
		&edgeNode{EdgeNameFront, idOffset + EdgeNameFront.Index(), EdgeNameFront.Dimension3D(), nil},
		&edgeNode{EdgeNameBack, idOffset + EdgeNameBack.Index(), EdgeNameBack.Dimension3D(), nil},
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

func New(paths []EdgePath, seed int64, chaos float64) *EdgePathAttributor {
	return &EdgePathAttributor{
		Paths:    paths,
		Seed:     seed,
		Chaos:    stepped.ValueToSteps(chaos, edgePathMaxChaos),
		MaxChaos: edgePathMaxChaos,
	}
}

func (attributor *EdgePathAttributor) AddAttributes(tree htree.Tree, attrs *attributors.NodeAttributer) error {
	rand.Seed(attributor.Seed)
	//epsilon := 0.0000001

	regionMap := htree.NewNodeRegionMap(tree, htree.Origin, htree.UnityScale)
	//dimensionLookup := htree.NewNodeDimensionMap(tree, htree.Origin, htree.UnityScale)
	matrix := algo.BuildAdjacencyMatrix(tree, regionMap)

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
	chaos := stepped.StepsToValue(attributor.Chaos, attributor.MaxChaos)

	for leafID, neighbors := range matrix {
		//fmt.Printf("leaf %d has %d neighbors\n", leafID, len(neighbors))

		for _, neighbor := range neighbors {
			if graph.HasEdgeBetween(simple.Node(leafID), simple.Node(neighbor.ID())) {
				continue
			}

			randWeight := 1 + chaos*rand.Float64()*float64(numLeaves)
			graph.SetWeightedEdge(&simple.WeightedEdge{F: simple.Node(leafID), T: simple.Node(neighbor.ID()), W: randWeight})
			//fmt.Printf("%d to %d = %f\n", leafID, neighbor.ID(), randWeight)

		}

		// Also need to check this neighbor against the edges
		dim := regionMap.Region(leafID).Dimension()

		for _, e := range edges {
			if e.dim.DistanceSquared(dim) < 0.0000001 {
				e.neighbors = append(e.neighbors, leafID)
			}
		}
	}

	for _, path := range attributor.Paths {
		fromNode := int64(-1)
		toNode := int64(-1)
		for _, edge := range edges {
			if edge.name == path.From {
				// Choose a random neighbor
				fromNode = int64(edge.neighbors[rand.Intn(len(edge.neighbors))])
			}
			if edge.name == path.To {
				toNode = int64(edge.neighbors[rand.Intn(len(edge.neighbors))])
			}
		}
		if fromNode < 0 || toNode < 0 {
			panic("path ends not found")
		}

		shortest := gpath.DijkstraFrom(simple.Node(fromNode), graph)
		p, w := shortest.To(simple.Node(toNode))
		fmt.Printf("path weight %f\n", w)

		for _, graphNode := range p {
			attrs.SetAttribute(htree.NodeID(graphNode.ID()), OnPathAttr, OnPathValue)
		}
	}

	return nil
}
