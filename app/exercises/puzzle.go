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
}

type PuzzleResult struct {
	IsCorrectAnswer bool
	IsFinish        bool
	Variants        []string
	AnsweredText    string
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
	}
}

func (p *Puzzle) Start() *PuzzleResult {
	return &PuzzleResult{
		IsCorrectAnswer: false,
		IsFinish:        false,
		Variants:        p.variants,
		AnsweredText:    "",
	}
}

func (p *Puzzle) HandleAnswer(answer string) *PuzzleResult {
	isCorrectAnswer := false

	if p.currentPosition < len(p.text) {
		isCorrectAnswer = strings.ToLower(p.text[p.currentPosition]) == answer

		if isCorrectAnswer {
			p.currentPosition += 1
			p.variants = computeVariants(p.text[p.currentPosition:], p.variants)
		}
	}

	return &PuzzleResult{
		IsCorrectAnswer: isCorrectAnswer,
		IsFinish:        p.currentPosition >= len(p.text),
		Variants:        p.variants,
		AnsweredText:    strings.Join(p.text[:p.currentPosition], " "),
	}
}
