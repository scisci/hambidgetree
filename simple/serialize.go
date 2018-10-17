package simple

import (
	"encoding/json"
	"errors"
	htree "github.com/scisci/hambidgetree"
	"sort"
)

const JSONVersion = 1

var InvalidSplitType = errors.New("Invalid split type")

// Serialize this like so
//
// ratios: (tree.ratios.ratios) // should marshal themselves with some kind of type flag, if strings, then expressions, otherwise floats
// xyRatioIndex:
// zyRatioIndex:
// root: id
// node: [id: id, split: hvd, index, left: id, right: id]

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type idSlice []htree.NodeID

func (p idSlice) Len() int           { return len(p) }
func (p idSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p idSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type jsonTree struct {
	Version      int          `json:"version"`
	Ratios       []string     `json:"ratios"`
	RatioIndexXY int          `json:"ratioIndexXY"`
	RatioIndexZY int          `json:"ratioIndexZY"`
	Root         htree.NodeID `json:"root"`
	Nodes        []jsonNode   `json:"nodes"`
}

type jsonBranch struct {
	SplitType  string       `json:"type"`
	LeftIndex  int          `json:"leftIndex"`
	RightIndex int          `json:"rightIndex"`
	Left       htree.NodeID `json:"left"`
	Right      htree.NodeID `json:"right"`
}

type jsonNode struct {
	ID     htree.NodeID `json:"id"`
	Branch *jsonBranch  `json:"branch,omitempty"`
}


func shortStringForSplitType(splitType htree.SplitType) string {
	switch splitType {
	case htree.SplitTypeHorizontal:
		return "h"
	case htree.SplitTypeVertical:
		return "v"
	case htree.SplitTypeDepth:
		return "d"
	default:
		return "-"
	}
}


func splitTypeForShortString(shortString string) (htree.SplitType, bool) {
	if shortString == "h" {
		return htree.SplitTypeHorizontal, true
	}

	if shortString == "v" {
		return htree.SplitTypeVertical, true
	}

	if shortString == "d" {
		return htree.SplitTypeDepth, true
	}

	return htree.SplitTypeHorizontal, false
}


func UnmarshalJSON(version int, data []byte) (*Tree, error) {
	tree := &Tree{}
	err := json.Unmarshal(data, &tree)
	return tree, err
}

func (tree *Tree) MarshalJSON() ([]byte, error) {
	// Build up a list of all nodes
	allIDs := make([]htree.NodeID, len(tree.nodes))
	index := 0
	for id, _ := range tree.nodes {
		allIDs[index] = id
		index++
	}
	sort.Sort(idSlice(allIDs))

	jNodes := make([]jsonNode, len(tree.nodes))
	for i, id := range allIDs {
		node, ok := tree.nodes[id]
		if !ok {
			panic("id error")
		}

		branch := node.branch
		var jBranch *jsonBranch
		if branch != nil {
			jBranch = &jsonBranch{
				SplitType:  shortStringForSplitType(branch.splitType),
				LeftIndex:  branch.leftIndex,
				RightIndex: branch.rightIndex,
				Left:       branch.left.id,
				Right:      branch.right.id,
			}
		}

		jNodes[i] = jsonNode{
			ID:     node.id,
			Branch: jBranch,
		}
	}

	jTree := jsonTree{
		Version:      0,
		Ratios:       tree.ratioSource.Exprs(),
		RatioIndexXY: tree.ratioIndexXY,
		RatioIndexZY: tree.ratioIndexZY,
		Root:         tree.root.ID(),
		Nodes:        jNodes,
	}

	return json.Marshal(jTree)
}

func (tree *Tree) UnmarshalJSON(data []byte) error {
	var jTree jsonTree
	if err := json.Unmarshal(data, &jTree); err != nil {
		return err
	}

	parents := make(map[htree.NodeID]htree.NodeID)
	nodes := make(map[htree.NodeID]*Node)

	for _, jNode := range jTree.Nodes {
		nodes[jNode.ID] = &Node{
			id: jNode.ID,
		}
	}

	for _, jNode := range jTree.Nodes {
		if jNode.Branch != nil {
			jBranch := jNode.Branch

			splitType, ok := splitTypeForShortString(jBranch.SplitType)
			if !ok {
				return InvalidSplitType
			}

			branch := NewBranch(
				splitType,
				nodes[jBranch.Left],
				nodes[jBranch.Right],
				jBranch.LeftIndex,
				jBranch.RightIndex,
			)

			nodes[jNode.ID].branch = branch
			parents[jBranch.Left] = jNode.ID
			parents[jBranch.Right] = jNode.ID
		}
	}

	ratioSource, err := htree.NewExprRatioSource(jTree.Ratios)
	if err != nil {
		return err
	}

	tree.nodes = nodes
	tree.parents = parents
	tree.ratioIndexXY = jTree.RatioIndexXY
	tree.ratioIndexZY = jTree.RatioIndexZY
	tree.ratioSource = ratioSource
	tree.root = nodes[jTree.Root]
	return nil
}
