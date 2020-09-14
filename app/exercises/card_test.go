package exercises

import (
	"testing"

	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look to it", RussianText: "Смотри на это"},
		{EnglishText: "I get it", RussianText: "Я понял это"},
		{EnglishText: "Check it out!", RussianText: "Зацени!"},
	}

	cards := NewCard(phrases, false)

	assert.Len(t, cards.queue, 6)
	assert.True(t, cards.queue[0].IsEnglishText)
	assert.Equal(t, "look to it", cards.queue[0].Phrase.EnglishText)
	assert.False(t, cards.queue[1].IsEnglishText)
	assert.Equal(t, "look to it", cards.queue[1].Phrase.EnglishText)
}

func TestKnowTrueNext(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look to it", RussianText: "Смотри на это"},
		{EnglishText: "I get it", RussianText: "Я понял это"},
	}

	cards := NewCard(phrases, false)
	assert.Len(t, cards.queue, 4)

	next := cards.Start()
	assert.False(t, next.IsFinish)
	assert.Equal(t, "look to it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "look to it", next.Card.Phrase.EnglishText)
	assert.False(t, next.Card.IsEnglishText)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.False(t, next.Card.IsEnglishText)

	next = cards.Next(true)
	assert.True(t, next.IsFinish)
	assert.Nil(t, next.Card)

	next = cards.Next(true)
	assert.True(t, next.IsFinish)
	assert.Nil(t, next.Card)
}

func TestLearnNext(t *testing.T) {
	phrases := []phrase.Phrase{
		{EnglishText: "look to it", RussianText: "Смотри на это"},
		{EnglishText: "I get it", RussianText: "Я понял это"},
	}

	cards := NewCard(phrases, false)
	assert.Len(t, cards.queue, 4)

	next := cards.Start()
	assert.False(t, next.IsFinish)
	assert.Equal(t, "look to it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 3)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "look to it", next.Card.Phrase.EnglishText)
	assert.False(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 2)

	next = cards.Next(false)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 2)
	// returns to the queue
	assert.Equal(t, cards.queue[1].Phrase.EnglishText, "look to it")
	assert.False(t, cards.queue[1].IsEnglishText)

	next = cards.Next(false)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.False(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 2)
	// returns to the queue
	assert.Equal(t, cards.queue[1].Phrase.EnglishText, "I get it")
	assert.True(t, cards.queue[1].IsEnglishText)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "look to it", next.Card.Phrase.EnglishText)
	assert.False(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 1)

	next = cards.Next(true)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 0)

	next = cards.Next(false)
	assert.False(t, next.IsFinish)
	assert.Equal(t, "I get it", next.Card.Phrase.EnglishText)
	assert.True(t, next.Card.IsEnglishText)
	assert.Len(t, cards.queue, 0)

	next = cards.Next(true)
	assert.True(t, next.IsFinish)
}
