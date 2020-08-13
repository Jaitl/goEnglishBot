package exercises

import (
	"math/rand"
	"time"

	"github.com/jaitl/goEnglishBot/app/phrase"
)

type Composite struct {
	mode         Mode
	phrases      []phrase.Phrase
	curExercises Exercise
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

	ex := NewExercise(mode, phrases[0].EnglishText)

	return &Composite{
		mode:         mode,
		phrases:      phrases,
		curExercises: *ex,
		curPos:       0,
		isFinish:     false,
	}
}

func (c *Composite) Next() *CompositePuzzleResult {
	if c.curExercises.IsFinish() && !c.isFinish {
		c.curExercises = *NewExercise(c.mode, c.phrases[c.curPos].EnglishText)
	}

	result := c.curExercises.Start()

	return &CompositePuzzleResult{
		IsFinish:     c.isFinish,
		Result:       result,
		Phrase:       &c.phrases[c.curPos],
		Pos:          c.curPos,
		CountPhrases: len(c.phrases),
	}
}

func (c *Composite) Skip() *CompositePuzzleResult {
	if !c.curExercises.IsFinish() {
		c.nextPos()
	}

	if !c.isFinish {
		c.curExercises = *NewExercise(c.mode, c.phrases[c.curPos].EnglishText)
	}

	result := c.curExercises.Start()

	return &CompositePuzzleResult{
		IsFinish:     c.isFinish,
		Result:       result,
		Phrase:       &c.phrases[c.curPos],
		Pos:          c.curPos,
		CountPhrases: len(c.phrases),
	}
}

func (c *Composite) HandleAnswer(answ string) *CompositePuzzleResult {
	pos := c.curPos
	result := c.curExercises.HandleAnswer(answ)

	if result.IsFinish {
		c.nextPos()
	}

	return &CompositePuzzleResult{
		IsFinish:     c.isFinish,
		Result:       result,
		Phrase:       &c.phrases[pos],
		Pos:          pos,
		CountPhrases: len(c.phrases),
	}
}

func (c *Composite) nextPos() {
	if c.isFinish {
		return
	}

	if c.curPos+1 >= len(c.phrases) {
		c.isFinish = true
	} else {
		c.curPos += 1
	}
}
