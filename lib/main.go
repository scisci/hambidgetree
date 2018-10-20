package main

/*
#include "stdlib.h" // for C.free

struct htree_vec {
	double x;
	double y;
	double z;
};

struct htree_region {
	long long id;
	struct htree_vec min;
	struct htree_vec max;
};

*/
import "C"

import (
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/factory"
	"github.com/scisci/hambidgetree/generators/randombasic"
	"github.com/scisci/hambidgetree/golden"
	"sync"
	"unsafe"
)

/*
//export FillRect
func FillRect() C.struct_HTreeRect {
	count := 2
	list := C.malloc(C.size_t(C.sizeof_double * count))
	for i := 0; i < count; i++ {
		ptr := unsafe.Pointer(uintptr(list) + uintptr(C.sizeof_double*i))
		*(*C.double)(ptr) = C.double(i)
	}
	return C.struct_HTreeRect{x: 1.0, y: 2.0, w: 3.0, h: 5.0, more: (*C.double)(list)}
}

//export FreeRect
func FreeRect(r C.struct_HTreeRect) {
	C.free(unsafe.Pointer(r.more))
}
*/

var treeRegMu sync.Mutex
var treeRegIndex = 0
var treeReg = make(map[int]htree.Tree)

//export CreateRandomBasicTree
func CreateRandomBasicTree(containerRatio float64, numRects int, seed int64) int {
	ratioSource := golden.RatioSource()
	gen, err := randombasic.New(ratioSource, containerRatio, numRects, seed)
	if err != nil {
		fmt.Printf("Failed to create generator (%v)\n", err)
		return -1
	}
	tree, err := gen.Generate()
	if err != nil {
		fmt.Printf("Failed to generate tree (%v)\n", err)
		return -1
	}
	return registerTree(tree)
}

//export CreateLeafRegions
func CreateLeafRegions(id int) (data_ptr unsafe.Pointer, count C.int) {
	tree := lookupTree(id)
	if tree == nil {
		fmt.Printf("Tree does not exist")
		return nil, 0
	}

	var leaves []htree.NodeRegion
	it := htree.NewRegionIterator(tree, htree.Origin, htree.UnityScale)
	for it.HasNext() {
		nodeRegion := it.Next()
		if nodeRegion.Node().Branch() == nil {
			// Its a leaf
			leaves = append(leaves, nodeRegion)
		}
	}

	// Allocate space
	count = C.int(len(leaves))
	data_ptr = C.malloc(C.size_t(C.sizeof_struct_htree_region * count))
	for i, region := range leaves {
		dim := region.Region().AlignedBox()
		c_node_region := C.struct_htree_region{}
		c_node_region.id = C.longlong(region.Node().ID())
		c_node_region.min = C.struct_htree_vec{C.double(dim.Left()), C.double(dim.Top()), C.double(dim.Front())}
		c_node_region.max = C.struct_htree_vec{C.double(dim.Right()), C.double(dim.Bottom()), C.double(dim.Back())}

		// Store the value in there
		ptr := unsafe.Pointer(uintptr(data_ptr) + uintptr(C.sizeof_struct_htree_region*i))
		*(*C.struct_htree_region)(ptr) = c_node_region
	}

	return
}

//export CreateSerializationOfTree
func CreateSerializationOfTree(id int) (data_ptr unsafe.Pointer, count C.int) {
	tree := lookupTree(id)
	data, err := factory.MarshalJSON(tree)
	if err != nil {
		fmt.Printf("Failed to marshal tree (%v)\n", err)
		return nil, 0
	}
	return C.CBytes(data), C.int(len(data))
}

//export Release
func Release(data_ptr unsafe.Pointer) {
	C.free(data_ptr)
}

//export CreateDeserializedTree
func CreateDeserializedTree(data_ptr unsafe.Pointer, count C.int) int {
	data := C.GoBytes(data_ptr, count)
	tree, err := factory.UnmarshalJSON(data)
	if err != nil {
		fmt.Printf("Failed to unmarshal tree (%v)\n", err)
		return -1
	}
	return registerTree(tree)
}

//export ReleaseTree
func ReleaseTree(id int) {
	treeRegMu.Lock()
	defer treeRegMu.Unlock()
	delete(treeReg, id)
	fmt.Printf("Released tree handle %d\n", id)
}

func registerTree(tree htree.Tree) int {
	treeRegMu.Lock()
	defer treeRegMu.Unlock()
	treeRegIndex++
	treeReg[treeRegIndex] = tree
	fmt.Printf("Registered tree handle %d\n", treeRegIndex)
	return treeRegIndex
}

func lookupTree(index int) htree.Tree {
	treeRegMu.Lock()
	defer treeRegMu.Unlock()
	return treeReg[index]
}

func main() {}
