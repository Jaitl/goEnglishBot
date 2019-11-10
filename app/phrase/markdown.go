package phrase

import (
	"strings"
)

const (
	rowPattern = "#_%v_ \"*%v*\": _%v_"
)

func ToMarkdownTable(ph []Phrase) string {

	if len(ph) == 0 {
		return "Список фраз пуст"
	}

	rows := make([]string, 0, len(ph))

	for _, p := range ph {
		rows = append(rows, p.ToMarkdown())
	}

	return strings.Join(rows, "\n")
}
