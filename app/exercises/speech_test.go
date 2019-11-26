package exercises

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpeechHandleCorrectAnswer(t *testing.T) {
	speech := NewSpeech("This tests whether a pattern-matches a string.")

	res := speech.HandleAnswer("This tests whether a")
	assert.False(t, res.IsCorrectAnswer)
	assert.Equal(t, float32(0.5), res.MatchScore)

	res = speech.HandleAnswer("This tests")
	assert.False(t, res.IsCorrectAnswer)
	assert.Equal(t, float32(0.25), res.MatchScore)

	res = speech.HandleAnswer("This tests whether a pattern ttt a wrong")
	assert.True(t, res.IsCorrectAnswer)
	assert.Equal(t, float32(0.75), res.MatchScore)
}
