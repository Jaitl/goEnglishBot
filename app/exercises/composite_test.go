package exercises

import (
	"testing"

	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/stretchr/testify/assert"
)

func TestCompositePuzzle(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check it"},
	}
	composite := NewComposite(phrases, PuzzleMode, false)

	assert.Len(t, composite.phrases, 3)
	assert.Equal(t, 0, composite.curPos)

	// 0
	result := composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)

	result = composite.HandleAnswer("look")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer("it")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// 1
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)
	assert.Equal(t, 1, composite.curPos)

	result = composite.HandleAnswer("i")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer("get")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// 2
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 2)
	assert.False(t, result.IsFinish)
	assert.Equal(t, 2, composite.curPos)

	result = composite.HandleAnswer("check")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 1)

	result = composite.HandleAnswer("it")
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)
	assert.Len(t, result.Result.Variants, 0)

	// correct
	result = composite.Next()
	assert.Len(t, result.Result.Variants, 0)
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.Equal(t, 2, composite.curPos)
}

func TestCompositeWrite(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check it"},
	}
	composite := NewComposite(phrases, WriteMode, false)

	assert.Len(t, composite.phrases, 3)
	assert.Equal(t, 0, composite.curPos)

	// 0
	result := composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 0, composite.curPos)

	result = composite.HandleAnswer("look it")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 1
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 1, composite.curPos)

	result = composite.HandleAnswer("i get")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 2
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 2, composite.curPos)

	result = composite.HandleAnswer("check it")
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// correct
	result = composite.Next()
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.Equal(t, 2, composite.curPos)
}

func TestCompositeSpeech(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check - it"},
	}

	composite := NewComposite(phrases, SpeechMode, false)

	assert.Len(t, composite.phrases, 3)
	assert.Equal(t, 0, composite.curPos)

	// 0
	result := composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 0, composite.curPos)

	result = composite.HandleAnswer("look it")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 1
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 1, composite.curPos)

	result = composite.HandleAnswer("i get")
	assert.False(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// 2
	result = composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, 2, composite.curPos)

	result = composite.HandleAnswer("check")
	assert.False(t, result.IsFinish)
	assert.False(t, result.Result.IsFinish)
	assert.False(t, result.Result.IsCorrectAnswer)

	result = composite.HandleAnswer("check it")
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.True(t, result.Result.IsCorrectAnswer)

	// correct
	result = composite.Next()
	assert.True(t, result.IsFinish)
	assert.True(t, result.Result.IsFinish)
	assert.Equal(t, 2, composite.curPos)
}

func TestCompositeSkip(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look it"},
		{EnglishText: "I get"},
		{EnglishText: "Check - it"},
	}

	composite := NewComposite(phrases, SpeechMode, false)

	result := composite.Next()
	assert.False(t, result.IsFinish)
	assert.Equal(t, "look it", result.Phrase.EnglishText)

	result = composite.Skip()
	assert.False(t, result.IsFinish)
	assert.Equal(t, "I get", result.Phrase.EnglishText)

	result = composite.Skip()
	assert.False(t, result.IsFinish)
	assert.Equal(t, "Check - it", result.Phrase.EnglishText)

	result = composite.Skip()
	assert.True(t, result.IsFinish)
}
