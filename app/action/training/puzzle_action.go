package training

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
}

const (
	Start          action.Stage = "start"          // Запускает puzzle тренировку
	WaitPushButton action.Stage = "waitPushButton" // Ожидает, когда пользователь выберет фразу

	Mode          action.SessionKey = "mode"
	PuzzleSession action.SessionKey = "puzzleSession"

	addedMsg string = "Добавлена фраза #%v \"%v\" с переводом \"%v\""
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

	keyboard := createKeyboard(puzzleRes.Variants)
	msg := fmt.Sprintf("Соберите фразу: %s", phrs.RussianText)

	err = a.Bot.SendWithKeyboard(cmd.GetUserId(), msg, keyboard)

	ses := action.CreateSession(cmd.GetUserId(), action.Add, WaitPushButton)
	ses.AddData(Mode, mode)
	ses.AddData(PuzzleSession, puzzle)
	a.ActionSession.UpdateSession(ses)

	return err
}

func (a *Action) waitPushButton(cmd command.Command, session *action.Session) error {
	return nil
}

func createKeyboard(variants []string) map[telegram.ButtonValue]telegram.ButtonName {
	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)

	for _, val := range variants {
		buttons[telegram.ButtonValue(val)] = telegram.ButtonName(val)
	}

	return buttons
}
