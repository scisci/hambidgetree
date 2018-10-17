package print

import (
	"bytes"
	"regexp"
	htree "github.com/scisci/hambidgetree"
)

func PrintRatios(ratioSource htree.RatioSource) string {
	r := regexp.MustCompile(`SQRT\(([^\)]+)\)`)
	exprs := ratioSource.Exprs()
	n := len(exprs)

	buf := bytes.NewBuffer(nil)
	buf.WriteString("[")
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		expr := r.ReplaceAllString(exprs[i], "√$1")

		// Replace SQRT(x) with symbol

		buf.WriteString(expr)
	}
	buf.WriteString("]")
	return buf.String()
}
