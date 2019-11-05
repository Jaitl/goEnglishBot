package remove

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

func (a *Action) GetType() action.Type {
	return action.Remove
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Remove: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		removeCmd := cmd.(*command.RemoveCommand)

		delCount, err := a.PhraseModel.RemovePhrase(removeCmd.UserId, removeCmd.IncNumber)

		if err != nil {
			return err
		}

		if delCount > 0 {
			return a.Bot.Send(cmd.GetUserId(), fmt.Sprintf("Фраза с id: %d удалена", removeCmd.IncNumber))
		}

		return nil
	}

	return fmt.Errorf("stage %s not found in RemoveAction", stage)
}
