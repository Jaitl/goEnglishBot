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
	Add    Type = "add"
	List   Type = "list"
	Card   Type = "card"
	Voice  Type = "voice"
	Me     Type = "me"
	Remove Type = "remove"
	Puzzle Type = "puzzle"
	Write  Type = "write"
	Speech Type = "speech"
)
