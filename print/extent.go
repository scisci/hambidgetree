package print

import (
	"fmt"
	htree "github.com/scisci/hambidgetree"
)

func PrintExtent(extent htree.Extent) string {
	return fmt.Sprintf("Extent{%.2f, %.2f}", extent.Start(), extent.End())
}
