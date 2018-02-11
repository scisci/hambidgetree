package randombasic

import (
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/golden"
	"testing"
)

func TestGenerator2D(t *testing.T) {
	ratios := htree.NewRatios([]float64{0.5, 1.0})
	treeRatios := htree.NewTreeRatios(ratios, 0.0000001)
	numLeaves := 3
	gen := New(treeRatios, 1, numLeaves, 1)
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	leaves := htree.FindLeaves(tree)
	if len(leaves) != numLeaves {
		t.Errorf("Got %d leaves, expected %d", len(leaves), numLeaves)
	}

	fmt.Println("Test2D")
	it := htree.NewImmutableNodeIterator(tree.Root())
	for it.HasNext() {
		node := it.Next()
		fmt.Printf("Node{%d, %v}\n", node.ID(), node.Branch())
	}

	rit := htree.NewRegionIterator(tree, htree.Origin, htree.UnityScale)
	for rit.HasNext() {
		rn := rit.Next()
		fmt.Printf("Node %d: %v\n", rn.Node().ID(), rn.Dimension())
	}

	/*
		dit := htree.NewDimensionalIterator(tree, htree.Origin, 1.0)
		for dit.HasNext() {
			dn := dit.Next()
			fmt.Printf("Node %d: %v\n", dn.ID(), dn.Dimension)
		}
	*/
}

func TestGenerator3D(t *testing.T) {
	ratios := golden.Ratios()
	treeRatios := htree.NewTreeRatios(ratios, 0.0000001)
	numLeaves := 50
	gen := New3D(treeRatios, 1, 1, numLeaves, 1)
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	leaves := htree.FindLeaves(tree)
	if len(leaves) != numLeaves {
		t.Errorf("Got %d leaves, expected %d", len(leaves), numLeaves)
	}
	/*
		fmt.Println("Test3D")
		it := NewNodeIterator(tree.root)
		for it.HasNext() {
			node := it.Next()
			fmt.Printf("Node{%d, %v}\n", node.ID(), node.Split())
		}
	*/

	rit := htree.NewRegionIterator(tree, htree.Origin, htree.UnityScale)
	for rit.HasNext() {
		rn := rit.Next()

		if rn.Node().Branch() != nil { // Ignore non-leaves
			continue
		}

		dimension := rn.Dimension()
		xy := dimension.Width() / dimension.Height()
		zx := dimension.Depth() / dimension.Width()
		zy := dimension.Depth() / dimension.Height()

		if xy > 8 || xy < 1.0/8.0 {
			fmt.Printf("Invalid xy ratio %f\n", xy)
		}

		if zx > 8 || zx < 1.0/8.0 {
			fmt.Printf("Invalid zx ratio %f\n", zx)
		}

		if zy > 8 || zy < 1.0/8.0 {
			fmt.Printf("Invalid zy ratio %f\n", zy)
		}
	}
	/*
		dit := htree.NewDimensionalIterator(tree, htree.Origin, 1.0)
		for dit.HasNext() {
			dn := dit.Next()
			//fmt.Printf("Node %d: %v\n", dn.ID(), dn.Dimension)

			if !dn.IsLeaf() {
				continue
			}

			xy := dn.Width() / dn.Height()
			zx := dn.Depth() / dn.Width()
			zy := dn.Depth() / dn.Height()

			if xy > 8 || xy < 1.0/8.0 {
				fmt.Printf("Invalid xy ratio %f\n", xy)
			}

			if zx > 8 || zx < 1.0/8.0 {
				fmt.Printf("Invalid zx ratio %f\n", zx)
			}

			if zy > 8 || zy < 1.0/8.0 {
				fmt.Printf("Invalid zy ratio %f\n", zy)
			}
		}
	*/
}
