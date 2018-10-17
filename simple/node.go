package simple

import (
	htree "github.com/scisci/hambidgetree"
)

type ParentLookup map[htree.NodeID]htree.NodeID
type NodeLookup map[htree.NodeID]*Node

type Node struct {
	id     htree.NodeID
	branch *Branch
}

func NewNode(id htree.NodeID, branch *Branch) *Node {
	return &Node{
		id:     id,
		branch: branch,
	}
}

func (n *Node) ID() htree.NodeID {
	return n.id
}

func (n *Node) Branch() htree.Branch {
	if n.branch == nil {
		return nil
	}

	return n.branch
}

type Branch struct {
	splitType  htree.SplitType
	left       *Node
	right      *Node
	leftIndex  int
	rightIndex int
}

func NewBranch(splitType htree.SplitType, left, right *Node, leftIndex, rightIndex int) *Branch {
	return &Branch{
		splitType:  splitType,
		left:       left,
		right:      right,
		leftIndex:  leftIndex,
		rightIndex: rightIndex,
	}
}

func (b *Branch) SplitType() htree.SplitType {
	return b.splitType
}

func (b *Branch) Left() htree.Node {
	return b.left
}

func (b *Branch) Right() htree.Node {
	return b.right
}

func (b *Branch) LeftIndex() int {
	return b.leftIndex
}

func (b *Branch) RightIndex() int {
	return b.rightIndex
}

type Tree struct {
	nodes        NodeLookup
	parents      ParentLookup
	ratioIndexXY int
	ratioIndexZY int
	ratioSource  htree.RatioSource
	root         *Node
}

func NewTree(ratioSource htree.RatioSource, ratioIndexXY, ratioIndexZY int, root *Node, nodes NodeLookup, parents ParentLookup) *Tree {
	if root == nil {
		panic("Can't create tree with nil root")
	}

	return &Tree{
		nodes:        nodes,
		parents:      parents,
		ratioIndexXY: ratioIndexXY,
		ratioIndexZY: ratioIndexZY,
		ratioSource:  ratioSource,
		root:         root,
	}
}

func (tree *Tree) Parent(id htree.NodeID) htree.Node {
	if parentID, ok := tree.parents[id]; ok {
		return tree.nodes[parentID]
	}

	return nil
}

func (tree *Tree) Node(id htree.NodeID) htree.Node {
	return tree.nodes[id]
}

func (tree *Tree) RatioSource() htree.RatioSource {
	return tree.ratioSource
}

func (tree *Tree) Root() htree.Node {
	return tree.root
}

func (tree *Tree) RatioIndexXY() int {
	return tree.ratioIndexXY
}

func (tree *Tree) RatioIndexZY() int {
	return tree.ratioIndexZY
}
