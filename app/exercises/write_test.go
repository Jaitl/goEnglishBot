package exercises

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWrite(t *testing.T) {
	text := NewWrite("She's hiding, bec-ause she's Embarrassed!!!")

	expText := []string{"she's", "hiding", "bec-ause", "she's", "embarrassed"}

	assert.Equal(t, expText, text.text)
	assert.Equal(t, 0, text.currentPosition)
}

func TestWriteHandleCorrectAnswer(t *testing.T) {
	text := NewWrite("the decrease in general-purpose")
	result := text.HandleAnswer([]string{"the"})

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 1)
	assert.Equal(t, result.NextAnswer, "decrease")
	assert.Equal(t, 3, result.WordsLeft)

	result = text.HandleAnswer([]string{"decrease", "in"})

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "general-purpose", result.NextAnswer)
	assert.Equal(t, 1, result.WordsLeft)

	result = text.HandleAnswer([]string{"general-purpose"})

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in general-purpose", result.AnsweredText)
	assert.Equal(t, 4, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)
}

func TestWriteHandleNotConnectAnswer(t *testing.T) {
	text := NewWrite("the decrease in general-purpose")
	result := text.HandleAnswer([]string{"not"})

	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 0)
	assert.Equal(t, result.NextAnswer, "the")
	assert.Equal(t, 4, result.WordsLeft)

	result = text.HandleAnswer([]string{"the", "decrease", "in"})

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 3)
	assert.Equal(t, result.NextAnswer, "general-purpose")
	assert.Equal(t, 1, result.WordsLeft)

	result = text.HandleAnswer([]string{"not"})

	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 3)
	assert.Equal(t, result.NextAnswer, "general-purpose")
	assert.Equal(t, 1, result.WordsLeft)

	result = text.HandleAnswer([]string{"general-purpose", "test", "not"})

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in general-purpose", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 4)
	assert.Equal(t, result.NextAnswer, "")
	assert.Equal(t, 0, result.WordsLeft)
}

func TestWriteHandlePositionBounds(t *testing.T) {
	text := NewWrite("the decrease in")
	result := text.HandleAnswer([]string{"the", "decrease", "in"})

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)

	result = text.HandleAnswer([]string{"test"})

	assert.False(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)
}
