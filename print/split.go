package print

import (
	htree "github.com/scisci/hambidgetree"
	"strconv"
)

func PrintSplit(split htree.Split) string {
	str := "Split{"
	switch split.Type() {
	case htree.SplitTypeHorizontal:
		str = str + "h"
	case htree.SplitTypeVertical:
		str = str + "v"
	case htree.SplitTypeDepth:
		str = str + "d"
	default:
		str = str + "?"
	}

	str = str + "," + strconv.Itoa(split.LeftIndex())
	str = str + "," + strconv.Itoa(split.RightIndex())
	str = str + "}"
	return str
}
