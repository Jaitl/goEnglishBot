package category_crud

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"strings"
)

type Action struct {
	Bot           *telegram.Telegram
	CategoryModel *category.Model
}

const (
	Start action.Stage = "start"
)

func (a *Action) GetType() action.Type {
	return action.Category
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.AddCategory: true, command.ListCategories: true, command.SetCategory: true, command.RemoveCategory: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		switch cm := cmd.(type) {
		case *command.AddCategoryCommand:
			return a.addAction(cm)
		case *command.ListCategoriesCommand:
			return a.listAction(cm)
		case *command.SetCategoriesCommand:
			return a.setAction(cm)
		case *command.RemoveCategoryCommand:
			return a.removeAction(cm)
		}
	}
	return fmt.Errorf("stage %s not found in CategoryAction", stage)
}

func (a *Action) addAction(cmd *command.AddCategoryCommand) error {
	incNumber, err := a.CategoryModel.NextIncNumberCategory(cmd.UserId)

	if err != nil {
		return err
	}

	err = a.CategoryModel.DeactivateAllCategories(cmd.UserId)

	if err != nil {
		return err
	}

	cat, err := a.CategoryModel.CreateCategory(cmd.UserId, incNumber, cmd.Name)

	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Добавлен каталог: %s", cat.ToMarkdown())

	return a.Bot.SendMarkdown(cmd.UserId, msg)
}

func (a *Action) listAction(cmd *command.ListCategoriesCommand) error {
	list, err := a.CategoryModel.AllCategories(cmd.GetUserId())

	if err != nil {
		return err
	}

	if len(list) == 0 {
		return a.Bot.Send(cmd.GetUserId(), "Нет ни одного каталога")
	}

	marks := make([]string, 0, len(list))

	for _, cat := range list {
		marks = append(marks, cat.ToMarkdown())
	}

	msg := strings.Join(marks, "\n")

	return a.Bot.SendMarkdown(cmd.GetUserId(), msg)
}

func (a *Action) setAction(cmd *command.SetCategoriesCommand) error {
	cat, err := a.CategoryModel.FindCategoryByIncNumber(cmd.UserId, cmd.IncNumber)
	if err != nil {
		return err
	}

	err = a.CategoryModel.DeactivateAllCategories(cmd.UserId)
	if err != nil {
		return err
	}

	err = a.CategoryModel.ActivateCategory(cmd.UserId, cmd.IncNumber)
	if err != nil {
		return err
	}

	cat.Active = true
	msg := fmt.Sprintf("Активирован каталог: %s", cat.ToMarkdown())

	return a.Bot.SendMarkdown(cmd.UserId, msg)
}

func (a *Action) removeAction(cmd *command.RemoveCategoryCommand) error {
	delCount, err := a.CategoryModel.RemoveCategory(cmd.UserId, cmd.IncNumber)
	if err != nil {
		return err
	}

	if delCount > 0 {
		return a.Bot.Send(cmd.GetUserId(), fmt.Sprintf("Каталог с id: %d удален", cmd.IncNumber))
	}

	return nil
}
