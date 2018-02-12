package main

/*
#include "stdlib.h" // for C.free
struct HTreeRect {
	double x;
	double y;
	double w;
	double h;
	double *more;
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
var treeReg = make(map[int]htree.ImmutableTree)

//export CreateRandomBasicTree
func CreateRandomBasicTree(containerRatio float64, numRects int, seed int64) int {
	ratios := golden.Ratios()
	gen := randombasic.New(ratios, containerRatio, numRects, seed)
	tree, err := gen.Generate()
	if err != nil {
		fmt.Printf("Failed to generate tree (%v)\n", err)
		return -1
	}
	return registerTree(tree)
}

//export SerializeTree
func SerializeTree(id int) (data_ptr unsafe.Pointer, count C.int) {
	tree := lookupTree(id)
	data, err := factory.MarshalJSON(tree)
	if err != nil {
		fmt.Printf("Failed to marshal tree (%v)\n", err)
		return nil, 0
	}
	return C.CBytes(data), C.int(len(data))
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

func registerTree(tree htree.ImmutableTree) int {
	treeRegMu.Lock()
	defer treeRegMu.Unlock()
	treeRegIndex++
	treeReg[treeRegIndex] = tree
	fmt.Printf("Registered tree handle %d\n", treeRegIndex)
	return treeRegIndex
}

func lookupTree(index int) htree.ImmutableTree {
	treeRegMu.Lock()
	defer treeRegMu.Unlock()
	return treeReg[index]
}

func main() {}
