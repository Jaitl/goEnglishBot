package phrase_add

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	CategoryModel *category.Model
}

const (
	Start               action.Stage = "start"               // Получает фразу на перевод
	WaitConfirm         action.Stage = "waitConfirm"         // Ожидает подтверждение автоматического перевода
	WaitCustomTranslate action.Stage = "waitCustomTranslate" // Ожидает когда пользователь пришлет свой перевод

	userPhrase   action.SessionKey = "userPhrase"
	awsTranslate action.SessionKey = "awsTranslate"
)

const (
	addMessageTemplate = "Добавлена фраза: %s\nв каталог: %s"
)

func (a *Action) GetType() action.Type {
	return action.PhraseAdd
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Text: true}
	case WaitConfirm:
		return map[command.Type]bool{command.KeyboardCallback: true}
	case WaitCustomTranslate:
		return map[command.Type]bool{command.Text: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitConfirm:
		return a.waitConfirmStage(cmd, session)
	case WaitCustomTranslate:
		return a.waitCustomTranslateStage(cmd, session)
	}

	return fmt.Errorf("stage %s not found in PhraseAddAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	var text string

	switch mcmd := cmd.(type) {
	case *command.TextCommand:
		text = phrase.Clear(mcmd.Text)
	default:
		return errors.New("command does not belong to Start stage in PhraseAddAction")
	}

	trans, err := a.AwsSession.Translate(text)
	if err != nil {
		return err
	}

	ses := action.CreateSession(cmd.GetUserId(), action.PhraseAdd, WaitConfirm)

	ses.AddData(userPhrase, text)
	ses.AddData(awsTranslate, trans)

	a.ActionSession.UpdateSession(ses)

	keyboard := map[telegram.ButtonValue]telegram.ButtonName{
		"save":   "Сохранить",
		"custom": "Свой перевод",
		"cancel": "Отменить",
	}

	err = a.Bot.SendWithKeyboard(cmd.GetUserId(), trans, keyboard)

	if err != nil {
		return err
	}

	return nil
}

func (a *Action) waitConfirmStage(cmd command.Command, session *action.Session) error {
	confirm, ok := cmd.(*command.KeyboardCallbackCommand)

	if !ok {
		return errors.New("command does not belong to ConfirmStage stage in PhraseAddAction")
	}

	switch confirm.Data {
	case "cancel":
		a.ActionSession.ClearSession(cmd.GetUserId())
		return a.Bot.Send(cmd.GetUserId(), "Добавление фразы отменено")

	case "save":
		a.ActionSession.ClearSession(cmd.GetUserId())
		incNumber, err := a.CategoryModel.NextIncNumberPhrase(confirm.UserId)
		if err != nil {
			return err
		}
		ph, cat, err := a.CategoryModel.CreatePhrase(cmd.GetUserId(), incNumber, session.GetStringData(userPhrase), session.GetStringData(awsTranslate))
		if err != nil {
			return err
		}

		msg := fmt.Sprintf(addMessageTemplate, ph.ToMarkdown(), cat.ToMarkdown())
		err = a.Bot.SendMarkdown(cmd.GetUserId(), msg)

		return err

	case "custom":
		session.Stage = WaitCustomTranslate
		a.ActionSession.UpdateSession(session)
		return a.Bot.Send(cmd.GetUserId(), "Отправте свой перевод")
	}

	return nil
}

func (a *Action) waitCustomTranslateStage(cmd command.Command, session *action.Session) error {
	translate, ok := cmd.(*command.TextCommand)

	if !ok {
		return errors.New("command does not belong to WaitCustomTranslate stage in PhraseAddAction")
	}

	a.ActionSession.ClearSession(cmd.GetUserId())

	incNumber, err := a.CategoryModel.NextIncNumberPhrase(translate.UserId)
	if err != nil {
		return err
	}

	ph, cat, err := a.CategoryModel.CreatePhrase(cmd.GetUserId(), incNumber, session.GetStringData(userPhrase), translate.Text)

	if err != nil {
		return err
	}

	msg := fmt.Sprintf(addMessageTemplate, ph.ToMarkdown(), cat.ToMarkdown())
	err = a.Bot.SendMarkdown(cmd.GetUserId(), msg)

	return err
}
