package hambidgetree

import "testing"
import "fmt"

func TestGenerator2D(t *testing.T) {
	ratios := NewRatios([]float64{0.5, 1.0})
	treeRatios := NewTreeRatios(ratios, 0.0000001)
	numLeaves := 3
	gen := NewRandomBasicTreeGenerator(treeRatios, 1, numLeaves, 1)
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	leaves := tree.Leaves()
	if len(leaves) != numLeaves {
		t.Errorf("Got %d leaves, expected %d", len(leaves), numLeaves)
	}

	fmt.Println("Test2D")
	it := NewNodeIterator(tree.root)
	for it.HasNext() {
		node := it.Next()
		fmt.Printf("Node{%d, %v}\n", node.ID(), node.Split())
	}

	dit := NewDimensionalIterator(tree, &Vector{0, 0, 0}, 1.0)
	for dit.HasNext() {
		dn := dit.Next()
		fmt.Printf("Node %d: %v\n", dn.ID(), dn.Dimension)
	}
}

func TestGenerator3D(t *testing.T) {
	ratios := NewGoldenRatios()
	treeRatios := NewTreeRatios(ratios, 0.0000001)
	numLeaves := 50
	gen := NewRandomBasic3DTreeGenerator(treeRatios, 1, 1, numLeaves, 1)
	tree, err := gen.Generate()
	if err != nil {
		t.Errorf("Error generating tree %v", err)
	}

	leaves := tree.Leaves()
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
	dit := NewDimensionalIterator(tree, &Vector{0, 0, 0}, 1.0)
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
}
