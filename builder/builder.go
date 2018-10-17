package builder

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/simple"
)

type uniqueIDGen struct {
	id htree.NodeID
}

func (u *uniqueIDGen) Next() htree.NodeID {
	u.id = u.id + 1
	return u.id
}

type dBranch struct {
	id         htree.NodeID
	left       *dNode
	right      *dNode
	leftIndex  int
	rightIndex int
	splitType  htree.SplitType
}

type dNode struct {
	id           htree.NodeID
	dimension    *htree.Dimension
	ratioIndexXY int
	ratioIndexZY int
}

func (node *dNode) ID() htree.NodeID {
	return node.id
}

func (node *dNode) Dimension() *htree.Dimension {
	return node.dimension
}

func (node *dNode) RatioIndexXY() int {
	return node.ratioIndexXY
}

func (node *dNode) RatioIndexZY() int {
	return node.ratioIndexZY
}

type TreeBuilder struct {
	idgen   *uniqueIDGen
	offset  *htree.Vector
	scale   float64
	regions map[htree.NodeID]*dNode
	//parents    map[htree.NodeID]htree.NodeID
	ratioSource htree.RatioSource
	root        *dNode
	branches    []*dBranch
	leaves      []*dNode
}

func New2D(ratioSource htree.RatioSource, ratioIndex int) *TreeBuilder {
	return New(ratioSource, ratioIndex, htree.RatioIndexUndefined)
}

func New3D(ratioSource htree.RatioSource, xyRatioIndex int, zyRatioIndex int) *TreeBuilder {
	return New(ratioSource, xyRatioIndex, zyRatioIndex)
}

func New(ratioSource htree.RatioSource, ratioIndexXY int, ratioIndexZY int) *TreeBuilder {
	offset := htree.Origin
	scale := htree.UnityScale
	regions := make(map[htree.NodeID]*dNode)

	ratios := ratioSource.Ratios()
	ratioXY := ratios[ratioIndexXY]
	ratioZY := 0.0

	if htree.IsRatioIndexDefined(ratioIndexZY) {
		ratioZY = ratios[ratioIndexZY]
	}

	max := htree.NewVector(ratioXY*scale, 1*scale, ratioZY*scale)

	idgen := &uniqueIDGen{}

	root := &dNode{
		id: idgen.Next(),
		dimension: htree.NewDimension3DV(
			offset,
			offset.Add(max)),
		ratioIndexXY: ratioIndexXY,
		ratioIndexZY: ratioIndexZY,
	}

	regions[root.id] = root

	return &TreeBuilder{
		idgen:       idgen,
		offset:      offset,
		scale:       scale,
		regions:     regions,
		ratioSource: ratioSource,
		root:        root,
		leaves:      []*dNode{root},
	}
}

func (b *TreeBuilder) Leaves() []htree.Leaf {
	leaves := make([]htree.Leaf, len(b.leaves))
	for i, v := range b.leaves {
		leaves[i] = v
	}
	return leaves
}

func (b *TreeBuilder) Branch(leafID htree.NodeID, splitType htree.SplitType, leftIndex, rightIndex int) (left, right htree.Leaf) {
	// Replace that node with a new node
	index := -1
	for i := 0; i < len(b.leaves); i++ {
		if b.leaves[i].id == leafID {
			index = i
			break
		}
	}

	if index == -1 {
		panic("Leaf is not part of this builder")
	}

	ratios := b.ratioSource.Ratios()

	var leftRegion, rightRegion htree.Region
	leaf := b.leaves[index]
	switch splitType {
	case htree.SplitTypeHorizontal:
		leftRegion, rightRegion = htree.SplitRegionHorizontal(ratios, leaf, leftIndex, rightIndex)
	case htree.SplitTypeVertical:
		leftRegion, rightRegion = htree.SplitRegionVertical(ratios, leaf, leftIndex, rightIndex)
	case htree.SplitTypeDepth:
		leftRegion, rightRegion = htree.SplitRegionDepth(ratios, leaf, leftIndex, rightIndex)
	default:
		panic("Unknown split type")
	}

	// Create a new node by
	leftNode := &dNode{
		id:           b.idgen.Next(),
		dimension:    leftRegion.Dimension(),
		ratioIndexXY: leftRegion.RatioIndexXY(),
		ratioIndexZY: leftRegion.RatioIndexZY(),
	}

	rightNode := &dNode{
		id:           b.idgen.Next(),
		dimension:    rightRegion.Dimension(),
		ratioIndexXY: rightRegion.RatioIndexXY(),
		ratioIndexZY: rightRegion.RatioIndexZY(),
	}

	//b.parents[leftNode.id] = leafID
	//b.parents[rightNode.id] = leafID
	b.regions[leftNode.id] = leftNode
	b.regions[rightNode.id] = rightNode
	b.branches = append(b.branches, &dBranch{id: leafID, left: leftNode, right: rightNode, leftIndex: leftIndex, rightIndex: rightIndex, splitType: splitType})

	// Remove the parent from the leaves and add the two new leaves
	b.leaves = append(b.leaves[:index], append(b.leaves[index+1:], leftNode, rightNode)...)
	return leftNode, rightNode
}

func (b *TreeBuilder) Build() (*simple.Tree, htree.RegionMap) {
	simpleBranches := make(map[htree.NodeID]*simple.Branch)
	simpleNodes := make(map[htree.NodeID]*simple.Node)
	simpleParents := make(map[htree.NodeID]htree.NodeID)
	// Start at the leaves and work up
	for i := len(b.branches) - 1; i >= 0; i-- {
		branch := b.branches[i]

		leftBranch := simpleBranches[branch.left.id]
		rightBranch := simpleBranches[branch.right.id]

		leftNode := simple.NewNode(branch.left.id, leftBranch)
		rightNode := simple.NewNode(branch.right.id, rightBranch)

		simpleNodes[branch.left.id] = leftNode
		simpleNodes[branch.right.id] = rightNode

		simpleParents[branch.left.id] = branch.id
		simpleParents[branch.right.id] = branch.id

		simpleBranches[branch.id] = simple.NewBranch(
			branch.splitType,
			leftNode,
			rightNode,
			branch.leftIndex,
			branch.rightIndex,
		)
	}

	var regionMap = make(map[htree.NodeID]htree.Region)
	for k, v := range b.regions {
		regionMap[k] = v
	}

	rootBranch := simpleBranches[b.root.id]
	root := simple.NewNode(b.root.id, rootBranch)
	simpleNodes[root.ID()] = root
	return simple.NewTree(b.ratioSource, b.root.ratioIndexXY, b.root.ratioIndexZY, root, simpleNodes, simpleParents), regionMap
}
