package puzzle

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/exercises"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	PhraseModel   *phrase.Model
	Audio         *telegram.AudioService
}

const (
	Start          action.Stage = "start"          // Запускает puzzle тренировку
	WaitPushButton action.Stage = "waitPushButton" // Ожидает, когда пользователь выберет фразу
)

const (
	Mode    action.SessionKey = "mode"
	Session action.SessionKey = "puzzleSession"
)

const (
	AudioMode string = "AudioMode"
	TransMode string = "TransMode"
)

func (a *Action) GetType() action.Type {
	return action.Puzzle
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.PuzzleAudio: true, command.PuzzleTrans: true}
	case WaitPushButton:
		return map[command.Type]bool{command.KeyboardCallback: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitPushButton:
		return a.waitPushButton(cmd, session)
	}

	return fmt.Errorf("stage %s not found in PuzzleAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	var mode string
	var incNumber int

	switch mcmd := cmd.(type) {
	case *command.PuzzleAudioCommand:
		mode = AudioMode
		incNumber = mcmd.IncNumber
	case *command.PuzzleTransCommand:
		mode = TransMode
		incNumber = mcmd.IncNumber
	default:
		return errors.New("command does not belong to Start stage in PuzzleAction")
	}

	phrs, err := a.PhraseModel.FindPhraseByIncNumber(cmd.GetUserId(), incNumber)

	if err != nil {
		return err
	}

	puzzle := exercises.NewPuzzle(phrs.EnglishText)
	puzzleRes := puzzle.Start()

	var msg string

	if mode == AudioMode {
		msg = "Соберите фразу, которую вы слышите"

		err := a.Audio.SendAudio(phrs)
		if err != nil {
			return err
		}
	} else {
		msg = fmt.Sprintf("Соберите фразу: %s", phrs.RussianText)
	}

	keyboard := createKeyboard(puzzleRes.Variants)

	err = a.Bot.SendWithKeyboard(cmd.GetUserId(), msg, keyboard)

	ses := action.CreateSession(cmd.GetUserId(), action.Puzzle, WaitPushButton)
	ses.AddData(Mode, mode)
	ses.AddData(Session, puzzle)
	a.ActionSession.UpdateSession(ses)

	return err
}

func (a *Action) waitPushButton(cmd command.Command, session *action.Session) error {
	callback, ok := cmd.(*command.KeyboardCallbackCommand)

	if !ok {
		return errors.New("command does not belong to WaitPushButton stage in PuzzleAction")
	}

	puzzle := session.Data[Session].(*exercises.Puzzle)

	puzzleRes := puzzle.HandleAnswer(callback.Data)

	if puzzleRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		return a.Bot.Send(cmd.GetUserId(), fmt.Sprintf("Фраза: %s\nУпражнение успешно завершено!", puzzleRes.AnsweredText))
	}

	msg := fmt.Sprintf("Фраза: %s", puzzleRes.AnsweredText)
	keyboard := createKeyboard(puzzleRes.Variants)

	return a.Bot.SendWithKeyboard(cmd.GetUserId(), msg, keyboard)
}

func createKeyboard(variants []string) map[telegram.ButtonValue]telegram.ButtonName {
	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)

	for _, val := range variants {
		buttons[telegram.ButtonValue(val)] = telegram.ButtonName(val)
	}

	return buttons
}
