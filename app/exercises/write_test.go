package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWrite(t *testing.T) {
	text := NewWrite("She's - hiding, because she's Embarrassed!!!")

	expText := []string{"she's", "hiding", "because", "she's", "embarrassed"}

	assert.Equal(t, expText, text.text)
	assert.Equal(t, 0, text.currentPosition)
}

func TestNewWrite2(t *testing.T) {
	text := NewWrite("Get the fuck out of here")

	expText := []string{"get", "the", "fuck", "out", "of", "here"}

	assert.Equal(t, expText, text.text)
	assert.Equal(t, 0, text.currentPosition)
}

func TestWriteHandleCorrectAnswer(t *testing.T) {
	text := NewWrite("the decrease in general-purpose")
	result := text.HandleAnswer("the")

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 1)
	assert.Equal(t, result.NextAnswer, "decrease")
	assert.Equal(t, 4, result.WordsLeft)

	result = text.HandleAnswer("decrease in")

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "general", result.NextAnswer)
	assert.Equal(t, 2, result.WordsLeft)

	result = text.HandleAnswer("general purpose")

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in general purpose", result.AnsweredText)
	assert.Equal(t, 5, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)
}

func TestWriteHandleNotConnectAnswer(t *testing.T) {
	text := NewWrite("the decrease in general-purpose")
	result := text.HandleAnswer("not")

	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 0)
	assert.Equal(t, result.NextAnswer, "the")
	assert.Equal(t, 5, result.WordsLeft)

	result = text.HandleAnswer("the decrease in")

	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 3)
	assert.Equal(t, result.NextAnswer, "general")
	assert.Equal(t, 2, result.WordsLeft)

	result = text.HandleAnswer("not")

	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 3)
	assert.Equal(t, result.NextAnswer, "general")
	assert.Equal(t, 2, result.WordsLeft)

	result = text.HandleAnswer("general purpose test not")

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in general purpose", result.AnsweredText)
	assert.Equal(t, text.currentPosition, 5)
	assert.Equal(t, result.NextAnswer, "")
	assert.Equal(t, 0, result.WordsLeft)
}

func TestWriteHandlePositionBounds(t *testing.T) {
	text := NewWrite("the decrease in")
	result := text.HandleAnswer("the decrease in")

	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)

	result = text.HandleAnswer("test")

	assert.False(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.Equal(t, "the decrease in", result.AnsweredText)
	assert.Equal(t, 3, text.currentPosition)
	assert.Equal(t, "", result.NextAnswer)
	assert.Equal(t, 0, result.WordsLeft)
}
