package learn_cards

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/exercises"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/utils"
	"time"
)

type Action struct {
	Bot           *telegram.Telegram
	CategoryModel *category.Model
	AwsSession    *aws.Session
	Audio         *telegram.AudioService
	ActionSession *action.SessionModel
}

const (
	Start                action.Stage = "start"
	WaitKnowButton       action.Stage = "waitKnowButton"
	WaitLearnedPhrase    action.Stage = "waitLearnedPhrase"
	WaitCheckCorrectness action.Stage = "waitCheckCorrectness"
)

const (
	Session   action.SessionKey = "cardSession"
	StartTime action.SessionKey = "startTime"
)

const (
	userKnowsAnswer       string = "know"
	userDoesntKnowsAnswer string = "not_know"
)

const (
	userApproveAnswer       string = "approve"
	userDoesntApproveAnswer string = "doesnt_approve"
)

const (
	userLearnedPhrase string = "userLearnedPhrase"
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
	case WaitKnowButton:
		return map[command.Type]bool{command.KeyboardCallback: true}
	case WaitLearnedPhrase:
		return map[command.Type]bool{command.KeyboardCallback: true}
	case WaitCheckCorrectness:
		return map[command.Type]bool{command.KeyboardCallback: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitKnowButton:
		return a.knowButtonAnswerStage(cmd, session)
	case WaitLearnedPhrase:
		return a.learnedAnswerStage(cmd, session)
	case WaitCheckCorrectness:
		return a.checkCorrectnessStage(cmd, session)
	}

	return fmt.Errorf("stage %s not found in LearnCardsAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	cardCmd := cmd.(*command.LearnCardsCommand)

	phrs, err := a.CategoryModel.SmartFindByRange(cmd.GetUserId(), cardCmd.From, cardCmd.To)

	if err != nil {
		return err
	}

	if len(phrs) == 0 {
		return errors.New("range doesn't correct")
	}

	cards := exercises.NewCard(phrs, true)

	_, err = a.nextCards(cardCmd.UserId, cards, true)

	if err != nil {
		return err
	}

	ses := action.CreateSession(cmd.GetUserId(), action.LearnCards, WaitKnowButton)
	ses.AddData(Session, cards)
	ses.AddData(StartTime, time.Now())
	a.ActionSession.UpdateSession(ses)

	return nil
}

func (a *Action) knowButtonAnswerStage(cmd command.Command, session *action.Session) error {
	callback, ok := cmd.(*command.KeyboardCallbackCommand)

	if !ok {
		return errors.New("command does not belong to WaitKnowButton stage in LearnCardsAction")
	}

	card := session.Data[Session].(*exercises.Card)

	switch callback.Data {
	case userKnowsAnswer:
		if err := a.checkCardCorrectness(cmd.GetUserId(), card); err != nil {
			return nil
		}
		session.Stage = WaitCheckCorrectness
		a.ActionSession.UpdateSession(session)
	case userDoesntKnowsAnswer:
		if err := a.showCard(cmd.GetUserId(), card); err != nil {
			return err
		}
		session.Stage = WaitLearnedPhrase
		a.ActionSession.UpdateSession(session)
	default:
		return fmt.Errorf("unknown button case: %s", callback.Data)
	}

	return nil
}

func (a *Action) learnedAnswerStage(cmd command.Command, session *action.Session) error {
	session.Stage = WaitKnowButton
	a.ActionSession.UpdateSession(session)

	card := session.Data[Session].(*exercises.Card)

	finish, err := a.nextCards(cmd.GetUserId(), card, false)
	if err != nil {
		return err
	}
	if finish {
		return a.finish(cmd.GetUserId(), session)
	}

	return nil
}

func (a *Action) checkCorrectnessStage(cmd command.Command, session *action.Session) error {
	callback, ok := cmd.(*command.KeyboardCallbackCommand)

	if !ok {
		return errors.New("command does not belong to WaitKnowButton stage in LearnCardsAction")
	}

	card := session.Data[Session].(*exercises.Card)


	switch callback.Data {
	case userApproveAnswer:
		finish, err := a.nextCards(cmd.GetUserId(), card, true)
		if err != nil {
			return err
		}
		if finish {
			return a.finish(cmd.GetUserId(), session)
		}
		session.Stage = WaitKnowButton
		a.ActionSession.UpdateSession(session)
	case userDoesntApproveAnswer:
		if err := a.showCard(cmd.GetUserId(), card); err != nil {
			return err
		}
		session.Stage = WaitLearnedPhrase
		a.ActionSession.UpdateSession(session)
	}

	return nil
}

func (a *Action) finish(chatId int, session *action.Session) error {
	a.ActionSession.ClearSession(chatId)
	startTime := session.Data[StartTime].(time.Time)
	elapsed := time.Since(startTime)
	msg := fmt.Sprintf("\nУпражнение успешно завершено за: %s!", utils.DurationPretty(elapsed))
	return a.Bot.Send(chatId, msg)
}

func (a *Action) nextCards(chatId int, c *exercises.Card, know bool) (bool, error) {
	cardRes := c.Next(know)
	if cardRes.IsFinish {
		return true, nil
	}

	msg := "Знаешь перевод фразы "

	if cardRes.Card.IsEnglishText {
		msg += fmt.Sprintf("\"%s\"?", cardRes.Card.Phrase.EnglishText)
	} else {
		msg += fmt.Sprintf("\"%s\"?", cardRes.Card.Phrase.RussianText)
	}

	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)
	buttons[telegram.ButtonValue(userKnowsAnswer)] = "Знаю"
	buttons[telegram.ButtonValue(userDoesntKnowsAnswer)] = "Не знаю"

	return false, a.Bot.SendWithKeyboard(chatId, msg, buttons)
}

func (a *Action) showCard(chatId int, c *exercises.Card) error {
	phr := c.CurCard().Phrase

	err := a.Bot.SendMarkdown(chatId, phr.ToMarkdown())
	if err != nil {
		return err
	}

	err = a.Audio.SendAudio(&phr)
	if err != nil {
		return err
	}

	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)
	buttons[telegram.ButtonValue(userLearnedPhrase)] = "Запомнил"

	return a.Bot.SendWithKeyboard(chatId, "Запомнил?", buttons)
}

func (a *Action) checkCardCorrectness(chatId int, c *exercises.Card) error {
	phr := c.CurCard().Phrase

	msg := "Перевод: "
	if c.CurCard().IsEnglishText {
		msg += fmt.Sprintf("\"%s\".", phr.RussianText)
	} else {
		msg += fmt.Sprintf("\"%s\".", phr.EnglishText)
	}

	msg += "\nУверен что знаешь?"

	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)
	buttons[telegram.ButtonValue(userApproveAnswer)] = "Да"
	buttons[telegram.ButtonValue(userDoesntApproveAnswer)] = "Нет"

	return a.Bot.SendWithKeyboard(chatId, msg, buttons)
}
