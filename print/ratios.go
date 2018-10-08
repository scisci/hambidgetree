package print

import (
	"bytes"
	"regexp"
	htree "github.com/scisci/hambidgetree"
)

func PrintRatios(ratios htree.Ratios) string {
	r := regexp.MustCompile(`SQRT\(([^\)]+)\)`)
	n := ratios.Len()

	buf := bytes.NewBuffer(nil)
	buf.WriteString("[")
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		expr := r.ReplaceAllString(ratios.Expr(i), "√$1")

		// Replace SQRT(x) with symbol

		buf.WriteString(expr)
	}
	buf.WriteString("]")
	return buf.String()
}