package add

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
)

type Action struct {
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	PhraseModel   *phrase.Model
}

const (
	Start               action.Stage = "start"               // Получает фразу на перевод
	WaitConfirm         action.Stage = "waitConfirm"         // Ожидает подтверждение автоматического перевода
	WaitCustomTranslate action.Stage = "waitCustomTranslate" // Ожидает когда пользователь пришлет свой перевод

	userPhrase   action.SessionKey = "userPhrase"
	awsTranslate action.SessionKey = "awsTranslate"

	addedMsg string = "Добавлена фраза #%v \"%v\" с переводом \"%v\""
)

func (a *Action) GetType() action.Type {
	return action.Add
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Add: true, command.Text: true}
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

	return fmt.Errorf("stage %s not found in AddAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	var text string

	switch mcmd := cmd.(type) {
	case *command.AddCommand:
		text = mcmd.Text
	case *command.TextCommand:
		text = mcmd.Text
	default:
		return errors.New("command does not belong to Start stage in AddAction")
	}

	trans, err := a.AwsSession.Translate(text)
	if err != nil {
		return err
	}

	ses := action.CreateSession(cmd.GetUserId(), action.Add, WaitConfirm)

	ses.AddData(userPhrase, text)
	ses.AddData(awsTranslate, trans)

	err = a.ActionSession.UpdateSession(ses)

	if err != nil {
		return err
	}

	keyboard := map[telegram.ButtonValue]telegram.ButtonName{"save": "Сохранить", "custom": "Свой перевод"}

	err = a.Bot.SendWithKeyboard(cmd.GetUserId(), trans, keyboard)

	if err != nil {
		return err
	}

	return nil
}

func (a *Action) waitConfirmStage(cmd command.Command, session *action.Session) error {
	confirm, ok := cmd.(*command.KeyboardCallbackCommand)

	if !ok {
		return errors.New("command does not belong to ConfirmStage stage in AddAction")
	}

	switch confirm.Data {
	case "save":
		err := a.ActionSession.ClearSession(cmd.GetUserId())
		if err != nil {
			return err
		}
		incNumber, err := a.PhraseModel.NextIncNumber(confirm.UserId)
		if err != nil {
			return err
		}
		err = a.PhraseModel.CreatePhrase(cmd.GetUserId(), incNumber, session.Data[userPhrase], session.Data[awsTranslate])
		if err != nil {
			return err
		}

		msg := fmt.Sprintf(addedMsg, incNumber, session.Data[userPhrase], session.Data[awsTranslate])
		err = a.Bot.Send(cmd.GetUserId(), msg)

		return err

	case "custom":
		session.Stage = WaitCustomTranslate
		err := a.ActionSession.UpdateSession(session)
		if err != nil {
			return err
		}
		err = a.Bot.Send(cmd.GetUserId(), "Отправте свой перевод")
		return err
	}

	return nil
}

func (a *Action) waitCustomTranslateStage(cmd command.Command, session *action.Session) error {
	translate, ok := cmd.(*command.TextCommand)

	if !ok {
		return errors.New("command does not belong to WaitCustomTranslate stage in AddAction")
	}

	err := a.ActionSession.ClearSession(cmd.GetUserId())

	if err != nil {
		return err
	}

	incNumber, err := a.PhraseModel.NextIncNumber(translate.UserId)
	if err != nil {
		return err
	}

	err = a.PhraseModel.CreatePhrase(cmd.GetUserId(), incNumber, session.Data[userPhrase], translate.Text)

	if err != nil {
		return err
	}

	msg := fmt.Sprintf(addedMsg, incNumber, session.Data[userPhrase], translate.Text)
	err = a.Bot.Send(cmd.GetUserId(), msg)

	return err
}
