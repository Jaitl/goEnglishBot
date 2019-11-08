package exercises

import (
	"strings"
)

type Write struct {
	text            []string
	currentPosition int
}

type WriteResult struct {
	IsCorrectAnswer bool
	IsFinish        bool
	AnsweredText    string
	NextAnswer      string
	WordsLeft       int
}

func NewWrite(text string) *Write {
	cText := ClearText(text)
	textParts := strings.Split(strings.ToLower(cText), " ")

	return &Write{
		text:            textParts,
		currentPosition: 0,
	}
}

func (p *Write) Start() *WriteResult {
	return &WriteResult{
		IsCorrectAnswer: false,
		IsFinish:        false,
		AnsweredText:    "",
		NextAnswer:      p.text[p.currentPosition],
		WordsLeft:       len(p.text),
	}
}

func (p *Write) HandleAnswer(answer []string) *WriteResult {
	isCorrectAnswer := false
	nextAnswer := ""

	if p.currentPosition < len(p.text) {
		var pos int
		isCorrectAnswer, pos = p.checkAnswer(answer)

		if isCorrectAnswer {
			p.currentPosition = pos
		}
	}

	if p.currentPosition < len(p.text) {
		nextAnswer = p.text[p.currentPosition]
	}

	return &WriteResult{
		IsCorrectAnswer: isCorrectAnswer,
		IsFinish:        p.currentPosition >= len(p.text),
		AnsweredText:    strings.Join(p.text[:p.currentPosition], " "),
		NextAnswer:      nextAnswer,
		WordsLeft:       len(p.text) - p.currentPosition,
	}
}

func (p *Write) checkAnswer(answer []string) (bool, int) {
	isCorrectAnswer := false
	pos := p.currentPosition
	for _, answ := range answer {
		if answ == p.text[pos] {
			isCorrectAnswer = true
			pos += 1
			if pos >= len(p.text) {
				break
			}
		}
	}

	return isCorrectAnswer, pos
}
