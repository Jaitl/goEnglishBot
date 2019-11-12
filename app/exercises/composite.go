package exercises

import (
	"github.com/jaitl/goEnglishBot/app/phrase"
	"math/rand"
	"time"
)

type Mode string

const (
	PuzzleMode Mode = "puzzle"
	WriteMode  Mode = "write"
)

type Composite struct {
	mode         Mode
	phrases      []phrase.Phrase
	curExercises interface{}
	curPos       int
	isFinish     bool
}

type CompositePuzzleResult struct {
	IsFinish     bool
	Result       *ExResult
	Phrase       *phrase.Phrase
	Pos          int
	CountPhrases int
}

func NewComposite(phrases []phrase.Phrase, mode Mode, random bool) *Composite {
	if random {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(phrases), func(i, j int) { phrases[i], phrases[j] = phrases[j], phrases[i] })
	}

	var ex interface{}

	switch mode {
	case PuzzleMode:
		ex = NewPuzzle(phrases[0].EnglishText)
	case WriteMode:
		ex = NewWrite(phrases[0].EnglishText)
	}
	return &Composite{
		mode:         mode,
		phrases:      phrases,
		curExercises: ex,
		curPos:       0,
		isFinish:     false,
	}
}

func (c *Composite) Next() *CompositePuzzleResult {
	var result *ExResult

	switch c.mode {
	case PuzzleMode:
		curEx := c.curExercises.(*Puzzle)
		if curEx.isFinish && !c.isFinish {
			curEx = NewPuzzle(c.phrases[c.curPos].EnglishText)
			c.curExercises = curEx
		}
		result = curEx.Start()
	case WriteMode:
		curEx := c.curExercises.(*Write)
		if curEx.isFinish && !c.isFinish {
			curEx = NewWrite(c.phrases[c.curPos].EnglishText)
			c.curExercises = curEx
		}
		result = curEx.Start()
	}

	return &CompositePuzzleResult{
		IsFinish:     c.isFinish,
		Result:       result,
		Phrase:       &c.phrases[c.curPos],
		Pos:          c.curPos,
		CountPhrases: len(c.phrases),
	}
}

func (c *Composite) HandleAnswer(answ []string) *CompositePuzzleResult {
	var result *ExResult
	pos := c.curPos

	switch c.mode {
	case PuzzleMode:
		result = c.curExercises.(*Puzzle).HandleAnswer(answ[0])
	case WriteMode:
		result = c.curExercises.(*Write).HandleAnswer(answ)
	}

	if result.IsFinish {
		if c.curPos+1 >= len(c.phrases) {
			c.isFinish = true
		} else {
			c.curPos += 1
		}
	}

	return &CompositePuzzleResult{
		IsFinish:     c.isFinish,
		Result:       result,
		Phrase:       &c.phrases[pos],
		Pos:          pos,
		CountPhrases: len(c.phrases),
	}
}
