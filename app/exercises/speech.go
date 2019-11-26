package exercises

import (
	"strings"
)

const (
	correctScore = 0.8
)

type Speech struct {
	text     map[string]int
	isFinish bool
}

func NewSpeech(text string) *Speech {
	index := createTextIndex(text)

	return &Speech{
		text:     index,
		isFinish: false,
	}
}

func (p *Speech) Start() *ExResult {
	return &ExResult{
		IsCorrectAnswer: false,
		IsFinish:        p.isFinish,
		MatchScore:      0,
	}
}

func (p *Speech) HandleAnswer(answer string) *ExResult {
	var score float32 = 0
	isCorrectAnswer := false

	if !p.isFinish {
		score = p.computeAnswerScore(answer)
		isCorrectAnswer = score >= correctScore

		if isCorrectAnswer {
			p.isFinish = true
		}
	}

	return &ExResult{
		IsCorrectAnswer: isCorrectAnswer,
		IsFinish:        p.isFinish,
		MatchScore:      score,
	}
}

func (p *Speech) computeAnswerScore(answer string) float32 {
	answerIndex := createTextIndex(answer)

	totalCount := 0
	totalAnswCount := 0

	for word, count := range p.text {
		totalCount += count

		answCount, ok := answerIndex[word]

		if ok {
			totalAnswCount += answCount
		}
	}

	return float32(totalAnswCount) / float32(totalCount)
}

func createTextIndex(test string) map[string]int {
	cText := ClearText(test)
	textParts := strings.Split(strings.ToLower(cText), " ")

	index := make(map[string]int)

	for _, part := range textParts {
		res, ok := index[part]
		if ok {
			index[part] = res + 1
		} else {
			index[part] = 1
		}
	}

	return index
}
