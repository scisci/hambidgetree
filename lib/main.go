package main

import "C"

import (
	"fmt"
	htree "github.com/scisci/hambidgetree"
	"github.com/scisci/hambidgetree/golden"
	"sync"
)

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
