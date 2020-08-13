package exercises

type Mode string

const (
	PuzzleMode Mode = "puzzle"
	WriteMode  Mode = "write"
	SpeechMode Mode = "speech"
)

type Exercise interface {
	Start() *ExResult
	HandleAnswer(answer string) *ExResult
	IsFinish() bool
}

type ExResult struct {
	IsCorrectAnswer bool
	IsFinish        bool
	Variants        []string
	AnsweredText    string
	NextAnswer      string
	WordsLeft       int
	MatchScore      float32
}

func NewExercise(mode Mode, phrase string) *Exercise {
	var ex Exercise

	switch mode {
	case PuzzleMode:
		ex = NewPuzzle(phrase)
	case WriteMode:
		ex = NewWrite(phrase)
	case SpeechMode:
		ex = NewSpeech(phrase)
	}

	return &ex
}
