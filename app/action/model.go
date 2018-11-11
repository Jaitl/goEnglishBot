package action

import (
	"github.com/jaitl/goEnglishBot/app/telegram/command"
)

type Action interface {
	GetType() Type
	GetStartStage() Stage
	GetStartCommands() []command.Type
	GetWaitCommands(stage Stage) map[command.Type]bool
	Execute(stage Stage, command command.Command, session *Session) error
}

type Type string
type Stage string

const (
	Add  Type = "add"
	List Type = "list"
)
