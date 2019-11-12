package phrase

import "strings"

func Clear(text string) string {
	str := strings.ReplaceAll(text, "‘", "'")
	str = strings.ReplaceAll(str, "`", "'")
	return str
}
