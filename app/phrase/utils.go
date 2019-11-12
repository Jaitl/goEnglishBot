package phrase

import (
	"golang.org/x/text/unicode/norm"
	"strings"
)

func Clear(text string) string {
	text = norm.NFKC.String(text)
	str := strings.ReplaceAll(text, "‘", "'")
	str = strings.ReplaceAll(str, "’", "'")
	str = strings.ReplaceAll(str, "`", "'")
	return str
}
