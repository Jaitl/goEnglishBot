package exercises

import (
	"math/rand"
	"strings"
	"time"
)

type Puzzle struct {
	text            []string
	variants        []string
	currentPosition int
	isFinish        bool
}

func NewPuzzle(text string) *Puzzle {
	cText := ClearText(text)
	textParts := strings.Split(strings.ToLower(cText), " ")
	uniqueParts := unique(textParts)
	variants := make([]string, len(uniqueParts))
	copy(variants, uniqueParts)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(variants), func(i, j int) { variants[i], variants[j] = variants[j], variants[i] })

	return &Puzzle{
		text:            textParts,
		variants:        variants,
		currentPosition: 0,
		isFinish:        false,
	}
}

func (p *Puzzle) Start() *ExResult {
	return &ExResult{
		IsCorrectAnswer: false,
		IsFinish:        p.isFinish,
		Variants:        p.variants,
		AnsweredText:    "",
		NextAnswer:      p.text[0],
		WordsLeft:       len(p.text),
	}
}

func (p *Puzzle) HandleAnswer(answer string) *ExResult {
	nextAnswer := ""
	isCorrectAnswer := false

	if !p.isFinish {
		isCorrectAnswer = strings.ToLower(p.text[p.currentPosition]) == answer

		if isCorrectAnswer {
			p.currentPosition += 1
			p.variants = computeVariants(p.text[p.currentPosition:], p.variants)
		}

		if p.currentPosition >= len(p.text) {
			p.isFinish = true
		} else {
			nextAnswer = p.text[p.currentPosition]
		}
	}

	return &ExResult{
		IsCorrectAnswer: isCorrectAnswer,
		IsFinish:        p.isFinish,
		Variants:        p.variants,
		AnsweredText:    strings.Join(p.text[:p.currentPosition], " "),
		NextAnswer:      nextAnswer,
		WordsLeft:       len(p.text) - p.currentPosition,
	}
}
