package print

import (
	"strconv"
	htree "github.com/scisci/hambidgetree"
)

func PrintComplements(c htree.Complements) string {
	str := ""
	for i := range c {
		str = str + "\n" + strconv.Itoa(i) + ":\n"
		for j := range c[i] {
			str = str + " " + PrintSplit(c[i][j])
		}
	}

	return str
}
