package me

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
	"strconv"
)

type Action struct {
	Bot *telegram.Telegram
}

const (
	Start action.Stage = "start"
)

func (a *Action) GetType() action.Type {
	return action.Me
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Me: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.Bot.Send(cmd.GetUserId(), "Ваш id: "+strconv.Itoa(cmd.GetUserId()))
	}

	return fmt.Errorf("stage %s not found in MeAction", stage)
}
