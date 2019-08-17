package list

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
)

type Action struct {
	Bot         *telegram.Telegram
	PhraseModel *phrase.Model
}

const (
	Start action.Stage = "start"
)

func (a *Action) GetType() action.Type {
	return action.List
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.List: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	}

	return fmt.Errorf("stage %s not found in ListAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	list, err := a.PhraseModel.AllPhrases(cmd.GetUserId())

	if err != nil {
		return err
	}

	message := phrase.ToMarkdownTable(list)

	err = a.Bot.SendMarkdown(cmd.GetUserId(), message)

	return err
}
