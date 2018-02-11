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
	goldenRatios := golden.Ratios()

	g2 := grid.New2D(2)
	g3 := grid.New3D(3)
	rb2, err := randombasic.New(goldenRatios, 1, 20, time.Now().UnixNano()).Generate()
	if err != nil {
		t.Errorf("failed to create random basic %v", err)
	}
	rb3, err := randombasic.New3D(goldenRatios, 1, 1, 20, time.Now().UnixNano()).Generate()
	if err != nil {
		t.Errorf("failed to create random basic %v", err)
	}

	treeTests := []htree.ImmutableTree{
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
