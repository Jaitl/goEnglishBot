package exercises

type ExResult struct {
	IsCorrectAnswer bool
	IsFinish        bool
	Variants        []string
	AnsweredText    string
	NextAnswer      string
	WordsLeft       int
	MatchScore      float32
}
