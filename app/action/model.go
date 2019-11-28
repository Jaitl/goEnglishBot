package action

import (
	"github.com/jaitl/goEnglishBot/app/command"
)

type Action interface {
	GetType() Type
	GetStartStage() Stage
	GetWaitCommands(stage Stage) map[command.Type]bool
	Execute(stage Stage, command command.Command, session *Session) error
}

type Type string
type Stage string

const (
	PhraseAdd    Type = "phraseAdd"
	PhrasesList  Type = "phrasesList"
	PhraseRemove Type = "phraseRemove"
	PhraseCard   Type = "phraseCard"
	Category     Type = "category"
	Voice        Type = "voice"
	Me           Type = "me"
	Puzzle       Type = "puzzle"
	Write        Type = "write"
	Speech       Type = "speech"
)
