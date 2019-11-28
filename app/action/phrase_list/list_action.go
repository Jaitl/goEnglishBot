package phrase_list

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	Bot           *telegram.Telegram
	CategoryModel *category.Model
}

const (
	Start action.Stage = "start"
)

const (
	phrasesInMessage int = 20
)

func (a *Action) GetType() action.Type {
	return action.PhrasesList
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.ListPhrases: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	}

	return fmt.Errorf("stage %s not found in PhrasesListAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	listCommand, ok := cmd.(*command.ListPhrasesCommand)

	if !ok {
		return fmt.Errorf("command %s not supported in PhrasesListAction", cmd.GetType())
	}

	var cat *category.Category
	var err error

	if listCommand.IncNumber == nil {
		cat, err = a.CategoryModel.FindActiveCategory(cmd.GetUserId())
	} else {
		cat, err = a.CategoryModel.FindCategoryByIncNumber(listCommand.GetUserId(), *listCommand.IncNumber)
	}

	if err != nil {
		return err
	}

	phraseCount := len(cat.Phrases)

	if phraseCount == 0 {
		return a.Bot.Send(cmd.GetUserId(), "Список фраз пуст")
	}

	messages := phrase.ToMarkdownTable(cat.Phrases, phrasesInMessage)

	err = a.Bot.SendMarkdown(cmd.GetUserId(), fmt.Sprintf("Каталог %s:", cat.ToMarkdown()))
	if err != nil {
		return err
	}

	for _, msg := range messages {
		err = a.Bot.SendMarkdown(cmd.GetUserId(), msg)
		if err != nil {
			return err
		}
	}

	return nil
}
