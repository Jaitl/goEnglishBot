package exercises

import (
	"math/rand"
	"regexp"
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
	cText := clearText(strings.ToLower(text))
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

func computeVariants(text []string, curVariants []string) []string {
	m := make(map[string]bool)
	variants := make([]string, 0, len(text))

	for _, val := range text {
		if _, ok := m[val]; !ok {
			m[val] = true
		}
	}

	for _, val := range curVariants {
		if _, ok := m[val]; ok {
			variants = append(variants, val)
		}
	}

	return variants
}

func clearText(text string) string {
	reg := regexp.MustCompile(`[^a-zA-Z1-9\s\\']+`)

	return reg.ReplaceAllString(text, "")
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
