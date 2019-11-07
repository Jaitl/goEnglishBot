package exercises

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPuzzle(t *testing.T) {
	puzzle := NewPuzzle("She's hiding, because she's Embarrassed!!!")

	expText := []string{"she's", "hiding", "because", "she's", "embarrassed"}
	assert.Equal(t, expText, puzzle.text)

	expVariants := []string{"she's", "hiding", "because", "embarrassed"}

	assert.Len(t, puzzle.variants, 4)
	assert.ElementsMatch(t, puzzle.variants, expVariants)

	assert.Equal(t, 0, puzzle.currentPosition)
}

func TestHandleCorrectAnswer(t *testing.T) {
	puzzle := NewPuzzle("She's hiding, because she's embarrassed 123!!!")

	result := puzzle.HandleAnswer("she's")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "hiding", "because", "embarrassed", "123"})
	assert.Equal(t, "she's", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 1)

	result = puzzle.HandleAnswer("hiding")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "because", "embarrassed", "123"})
	assert.Equal(t, "she's hiding", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 2)

	result = puzzle.HandleAnswer("because")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "embarrassed", "123"})
	assert.Equal(t, "she's hiding because", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 3)

	result = puzzle.HandleAnswer("she's")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"embarrassed", "123"})
	assert.Equal(t, "she's hiding because she's", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 4)

	result = puzzle.HandleAnswer("embarrassed")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"123"})
	assert.Equal(t, "she's hiding because she's embarrassed", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 5)

	result = puzzle.HandleAnswer("123")
	assert.True(t, result.IsCorrectAnswer)
	assert.True(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{})
	assert.Equal(t, "she's hiding because she's embarrassed 123", result.AnsweredText)
	assert.Equal(t, puzzle.currentPosition, 6)
}

func TestHandleNotConnectAnswer(t *testing.T) {
	puzzle := NewPuzzle("She's hiding, because she's embarrassed!!!")

	result := puzzle.HandleAnswer("test")
	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "hiding", "because", "embarrassed"})
	assert.Equal(t, puzzle.currentPosition, 0)

	result = puzzle.HandleAnswer("she's")
	assert.True(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "hiding", "because", "embarrassed"})
	assert.Equal(t, puzzle.currentPosition, 1)

	result = puzzle.HandleAnswer("not")
	assert.False(t, result.IsCorrectAnswer)
	assert.False(t, result.IsFinish)
	assert.ElementsMatch(t, result.Variants, []string{"she's", "hiding", "because", "embarrassed"})
	assert.Equal(t, puzzle.currentPosition, 1)
}
