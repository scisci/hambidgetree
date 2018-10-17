package print

import (
	"strconv"
	htree "github.com/scisci/hambidgetree"
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
	
		str = str + "," + strconv.Itoa(split.LeftIndex()-1)
		str = str + "," + strconv.Itoa(split.RightIndex()-1)
		str = str + "}"
		return str
}