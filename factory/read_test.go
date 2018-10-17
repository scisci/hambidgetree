package factory_test

import (
	"bytes"
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/factory"
	"github.com/scisci/hambidgetree/generators/grid"
	"github.com/scisci/hambidgetree/generators/randombasic"
	"github.com/scisci/hambidgetree/golden"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {
	goldenRatios := golden.RatioSource()

	g2 := grid.New2D(2)
	g3 := grid.New3D(3)
	gen2, err := randombasic.New(goldenRatios, 1, 20, time.Now().UnixNano())
	if err != nil {
		t.Errorf("failed to create generator %v", err)
	}
	rb2, err := gen2.Generate()
	if err != nil {
		t.Errorf("failed to create random basic %v", err)
	}
	gen3, err := randombasic.New3D(goldenRatios, 1, 1, 20, time.Now().UnixNano())
	if err != nil {
		t.Errorf("failed to create generator %v", err)
	}
	rb3, err := gen3.Generate()
	if err != nil {
		t.Errorf("failed to create random basic %v", err)
	}

	treeTests := []htree.Tree{
		g2, g3, rb2, rb3,
	}

	for i, tree := range treeTests {
		treeData, err := factory.MarshalJSON(tree)
		if err != nil {
			t.Errorf("test %d failed to marshal %v", i, err)
		}

		tree2, err := factory.UnmarshalJSON(treeData)
		if err != nil {
			t.Errorf("test %d to unmarshal %v", i, err)
		}

		treeData2, err := factory.MarshalJSON(tree2)
		if err != nil {
			t.Errorf("test %d failed to marshal 2nd time %v", i, err)
		}

		if bytes.Compare(treeData, treeData2) != 0 {
			fmt.Println(string(treeData))
			fmt.Println(string(treeData2))
			t.Errorf("test %d encoding/decoding non-symmetric", i)
		}
	}
}
