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
	"github.com/scisci/hambidgetree/golden"
	"sync"
	"unsafe"
)

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

var ratioRegistryMu sync.Mutex
var ratioRegistryIndex = 0
var ratioRegistry = make(map[int]htree.Ratios)

//export CreateGoldenRatios
func CreateGoldenRatios() int {
	return registerRatios(golden.Ratios())
}

//export ReleaseRatios
func ReleaseRatios(index int) {
	ratioRegistryMu.Lock()
	defer ratioRegistryMu.Unlock()
	delete(ratioRegistry, index)
}

//export PrintRatios
func PrintRatios(index int) {
	ratios := lookupRatios(index)
	fmt.Printf("%v", ratios)
}

func registerRatios(ratios htree.Ratios) int {
	ratioRegistryMu.Lock()
	defer ratioRegistryMu.Unlock()
	ratioRegistryIndex++
	ratioRegistry[ratioRegistryIndex] = ratios
	return ratioRegistryIndex
}

func lookupRatios(index int) htree.Ratios {
	ratioRegistryMu.Lock()
	defer ratioRegistryMu.Unlock()
	return ratioRegistry[index]
}

func main() {}
