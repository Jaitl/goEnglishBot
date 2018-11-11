package phrase

import (
	"fmt"
	"strings"
)

const (
	rowPattern = "#_%v_ \"*%v*\": _%v_"
)

func ToMarkdownTable(ph []Phrase) string {

	if len(ph) == 0 {
		return "Список фраз пуст"
	}

	rows := make([]string, len(ph))

	for _, p := range ph {
		rows = append(rows, prepareTableRow(&p))
	}

	return strings.Join(rows, "\n")
}

func prepareTableRow(p *Phrase) string {
	return fmt.Sprintf(rowPattern, p.IncNumber, p.EnglishText, p.RussianText)
}
