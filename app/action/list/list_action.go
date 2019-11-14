package list

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	Bot         *telegram.Telegram
	PhraseModel *phrase.Model
}

const (
	Start action.Stage = "start"
)

const (
	phrasesInMessage int = 20
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

	if len(list) == 0 {
		return a.Bot.Send(cmd.GetUserId(), "Список фраз пуст")
	}

	messages := phrase.ToMarkdownTable(list, phrasesInMessage)

	for _, msg := range messages {
		err = a.Bot.SendMarkdown(cmd.GetUserId(), msg)
		if err != nil {
			return err
		}
	}

	return a.Bot.Send(cmd.GetUserId(), fmt.Sprintf("Количество фраз: %d", len(list)))
}
