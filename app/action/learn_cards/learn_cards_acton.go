package learn_cards

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	Bot           *telegram.Telegram
	CategoryModel *category.Model
	AwsSession    *aws.Session
	Audio         *telegram.AudioService
}

const (
	Start action.Stage = "start"
)

func (a *Action) GetType() action.Type {
	return action.LearnCards
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.LearnCards: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd, session)
	}

	return fmt.Errorf("stage %s not found in LearnCardsAction", stage)
}

func (a *Action) startStage(cmd command.Command, session *action.Session) error {
	return a.Bot.Send(cmd.GetUserId(), "It is learn card action!")
}
