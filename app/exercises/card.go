package exercises

import (
	"math/rand"
	"time"

	"github.com/jaitl/goEnglishBot/app/phrase"
)

type CardPhrase struct {
	Phrase        phrase.Phrase
	IsEnglishText bool
}

type CardResult struct {
	Card     *CardPhrase
	IsFinish bool
}

type Card struct {
	queue []CardPhrase
	curCard *CardPhrase
}

func NewCard(phrases []phrase.Phrase, random bool) *Card {
	var q []CardPhrase

	for _, phrase := range phrases {
		q = append(q, CardPhrase{Phrase: phrase, IsEnglishText: true})
		q = append(q, CardPhrase{Phrase: phrase, IsEnglishText: false})
	}

	if random {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	}

	return &Card{
		queue: q,
		curCard: nil,
	}
}

func (c *Card) Start() *CardResult {
	return c.Next(true)
}

func (c *Card) Next(know bool) *CardResult {
	if !know && c.curCard != nil {
		c.queue = append(c.queue, *c.curCard)
	}

	if len(c.queue) > 0 {
		c.curCard = &c.queue[0]
		c.queue = c.queue[1:]
	} else {
		c.curCard = nil
	}

	return &CardResult{
		Card: c.curCard,
		IsFinish: c.curCard == nil,
	}
}
