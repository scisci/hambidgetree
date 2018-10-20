package grid

import (
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/builder"
	"github.com/scisci/hambidgetree/simple"
)

// Useful for testing, each level splits the tree in half, starting with a
// square
// Even levels split all leaves in half using a vertical line
// Odd levels split those splits in half with a horizontal line creating squares again
func New2D(levels int) *simple.Tree {
	ratioSource, err := htree.NewBasicRatioSource([]float64{0.5, 1.0})
	if err != nil {
		panic(err)
	}
	builder := builder.New2D(ratioSource, 1)

	for i := 0; i < levels; i++ {
		leaves := builder.Leaves()
		for _, leaf := range leaves {
			if i&1 == 0 {
				builder.Branch(leaf.ID(), htree.SplitTypeVertical, 0, 0)
			} else {
				builder.Branch(leaf.ID(), htree.SplitTypeHorizontal, 1, 1)
			}
		}
	}

	tree, _ := builder.Build()
	return tree
}

func New3D(levels int) *simple.Tree {
	ratioSource, err := htree.NewBasicRatioSource([]float64{0.5, 1.0, 2.0})
	if err != nil {
		panic(err)
	}
	builder := builder.New3D(ratioSource, 1, 1)

	for i := 0; i < levels; i++ {
		leaves := builder.Leaves()
		for _, leaf := range leaves {
			index := i % 3
			switch index {
			case 0:
				builder.Branch(leaf.ID(), htree.SplitTypeVertical, 0, 0)
			case 1:
				builder.Branch(leaf.ID(), htree.SplitTypeHorizontal, 1, 1)
			case 2:
				builder.Branch(leaf.ID(), htree.SplitTypeDepth, 1, 1)
			}
		}
	}

	tree, _ := builder.Build()
	return tree
}
