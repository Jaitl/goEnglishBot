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
	var from, to *int

	switch mcmd := cmd.(type) {
	case *command.PuzzleAudioCommand:
		mode = AudioMode
		from = mcmd.From
		to = mcmd.To
	case *command.PuzzleTransCommand:
		mode = TransMode
		from = mcmd.From
		to = mcmd.To
	default:
		return errors.New("command does not belong to Start stage in PuzzleAction")
	}

	phrs, err := a.PhraseModel.SmartFindByRange(cmd.GetUserId(), from, to)

	if err != nil {
		return err
	}

	if len(phrs) == 0 {
		return errors.New("don't correct range")
	}

	puzzle := exercises.NewComposite(phrs, exercises.PuzzleMode, true)
	err = a.newPhrase(puzzle, mode)

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

	puzzle := session.Data[Session].(*exercises.Composite)
	mode := session.GetStringData(Mode)

	puzzleRes := puzzle.HandleAnswer([]string{callback.Data})

	msg := fmt.Sprintf("Фраза №%d из %d", puzzleRes.Pos+1, puzzleRes.CountPhrases)
	msg += fmt.Sprintf("\nФраза: %s", puzzleRes.Result.AnsweredText)

	if puzzleRes.Result.IsFinish && puzzleRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		msg += "\nФраза успешно завершена!"
		msg += "\nУпражнение успешно завершено!"
		return a.Bot.Send(cmd.GetUserId(), msg)
	}

	if puzzleRes.Result.IsFinish && !puzzleRes.IsFinish {
		msg += "\nФраза успешно завершена!"
		err := a.Bot.Send(cmd.GetUserId(), msg)

		if err != nil {
			return err
		}

		return a.newPhrase(puzzle, mode)
	}

	keyboard := createKeyboard(puzzleRes.Result.Variants)
	return a.Bot.SendWithKeyboard(cmd.GetUserId(), msg, keyboard)
}

func (a *Action) newPhrase(puzzle *exercises.Composite, mode string) error {
	puzzleRes := puzzle.Next()

	msg := fmt.Sprintf("Фраза №%d из %d", puzzleRes.Pos+1, puzzleRes.CountPhrases)

	if mode == AudioMode {
		msg += "\nСоберите фразу, которую вы слышите"

		err := a.Audio.SendAudio(puzzleRes.Phrase)
		if err != nil {
			return err
		}
	} else {
		msg += fmt.Sprintf("\nСоберите фразу: %s", puzzleRes.Phrase.RussianText)
	}

	keyboard := createKeyboard(puzzleRes.Result.Variants)

	return a.Bot.SendWithKeyboard(puzzleRes.Phrase.UserId, msg, keyboard)
}

func createKeyboard(variants []string) map[telegram.ButtonValue]telegram.ButtonName {
	buttons := make(map[telegram.ButtonValue]telegram.ButtonName)

	for _, val := range variants {
		buttons[telegram.ButtonValue(val)] = telegram.ButtonName(val)
	}

	return buttons
}
