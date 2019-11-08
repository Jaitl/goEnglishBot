package write

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/exercises"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"strings"
)

type Action struct {
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	PhraseModel   *phrase.Model
	Audio         *telegram.AudioService
}

const (
	Start           action.Stage = "start"
	WaitWrittenText action.Stage = "waitWrittenText"
)

const (
	Mode        action.SessionKey = "mode"
	Session     action.SessionKey = "writeSession"
	CountErrors action.SessionKey = "countErrors"
)

const (
	AudioMode string = "AudioMode"
	TransMode string = "TransMode"
)

const (
	maxCountErrors int = 3
)

func (a *Action) GetType() action.Type {
	return action.Write
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.WriteAudio: true, command.WriteTrans: true}
	case WaitWrittenText:
		return map[command.Type]bool{command.Text: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitWrittenText:
		return a.waitWrittenText(cmd, session)
	}

	return fmt.Errorf("stage %s not found in WriteAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	var mode string
	var incNumber int

	switch mcmd := cmd.(type) {
	case *command.WriteAudioCommand:
		mode = AudioMode
		incNumber = mcmd.IncNumber
	case *command.WriteTransCommand:
		mode = TransMode
		incNumber = mcmd.IncNumber
	default:
		return errors.New("command does not belong to Start stage in WriteAction")
	}

	phrs, err := a.PhraseModel.FindPhraseByIncNumber(cmd.GetUserId(), incNumber)

	if err != nil {
		return err
	}

	write := exercises.NewWrite(phrs.EnglishText)

	var msg string

	if mode == AudioMode {
		msg = "Напишите фразу, которую вы слышите"

		err := a.Audio.SendAudio(phrs)
		if err != nil {
			return err
		}
	} else {
		msg = fmt.Sprintf("Напишите фразу: %s", phrs.RussianText)
	}

	err = a.Bot.Send(cmd.GetUserId(), msg)

	ses := action.CreateSession(cmd.GetUserId(), action.Write, WaitWrittenText)
	ses.AddData(Mode, mode)
	ses.AddData(Session, write)
	ses.AddData(CountErrors, 0)
	a.ActionSession.UpdateSession(ses)

	return err
}

func (a *Action) waitWrittenText(cmd command.Command, session *action.Session) error {
	text, ok := cmd.(*command.TextCommand)

	if !ok {
		return errors.New("command does not belong to WaitWrittenText stage in WriteAction")
	}

	write := session.Data[Session].(*exercises.Write)

	words := strings.Split(exercises.ClearText(text.Text), " ")

	writeRes := write.HandleAnswer(words)

	if writeRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		return a.Bot.Send(cmd.GetUserId(), fmt.Sprintf("Фраза: %s\nУпражнение успешно завершено!", writeRes.AnsweredText))
	}

	countErrors := session.GetIntData(CountErrors)

	msg := fmt.Sprintf("Фраза: %s\nОсталось слов: %d", writeRes.AnsweredText, writeRes.WordsLeft)

	if !writeRes.IsCorrectAnswer {
		msg += "\nНекорректное слово!"
		countErrors += 1
	}

	if countErrors >= maxCountErrors {
		msg += fmt.Sprintf("\nСледующее слово: %s", writeRes.NextAnswer)
		countErrors = 0
	}

	session.AddData(CountErrors, countErrors)
	a.ActionSession.UpdateSession(session)

	return a.Bot.Send(cmd.GetUserId(), msg)
}
