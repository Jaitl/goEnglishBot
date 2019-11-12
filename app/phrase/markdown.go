package phrase

import (
	"strings"
)

func ToMarkdownTable(ph []Phrase, inMessage int) []string {
	var rows [][]string
	row := make([]string, 0, inMessage)

	curCnt := 0
	for _, p := range ph {
		curCnt += 1
		row = append(row, p.ToMarkdown())
		if curCnt >= inMessage {
			curCnt = 0
			rows = append(rows, row)
			row = make([]string, 0, inMessage)
		}
	}

	if curCnt > 0 {
		rows = append(rows, row)
	}

	messages := make([]string, 0, len(rows))

	for _, r := range rows {
		msg := strings.Join(r, "\n")
		messages = append(messages, msg)
	}

	return messages
}
