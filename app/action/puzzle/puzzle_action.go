package puzzle

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
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	CategoryModel *category.Model
	Audio         *telegram.AudioService
}

const (
	Start          action.Stage = "start"          // Запускает puzzle тренировку
	WaitPushButton action.Stage = "waitPushButton" // Ожидает, когда пользователь выберет фразу
)

const (
	Mode         action.SessionKey = "mode"
	Session      action.SessionKey = "puzzleSession"
	ErrorsToHelp action.SessionKey = "errorsToHelp"
	CountErrors  action.SessionKey = "countErrors"
	StartTime    action.SessionKey = "startTime"
)

const (
	maxCountErrors int = 3
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

	phrs, err := a.CategoryModel.SmartFindByRange(cmd.GetUserId(), from, to)

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
	ses.AddData(ErrorsToHelp, 0)
	ses.AddData(CountErrors, 0)
	ses.AddData(StartTime, time.Now())
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
	countErrors := session.GetIntData(CountErrors)

	puzzleRes := puzzle.HandleAnswer(callback.Data)

	msg := fmt.Sprintf("Фраза №%d из %d", puzzleRes.Pos+1, puzzleRes.CountPhrases)
	msg += fmt.Sprintf("\nФраза: %s", puzzleRes.Result.AnsweredText)

	if puzzleRes.Result.IsFinish && puzzleRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		msg += fmt.Sprintf("\nПеревод: %s", puzzleRes.Phrase.RussianText)
		msg += "\nФраза успешно завершена!"
		msg += fmt.Sprintf("\nКоличество ошибок: %d", countErrors)

		startTime := session.Data[StartTime].(time.Time)
		elapsed := time.Since(startTime)

		msg += fmt.Sprintf("\nУпражнение успешно завершено за: %s!", utils.DurationPretty(elapsed))

		return a.Bot.Send(cmd.GetUserId(), msg)
	}

	if puzzleRes.Result.IsFinish && !puzzleRes.IsFinish {
		msg += fmt.Sprintf("\nПеревод: %s", puzzleRes.Phrase.RussianText)
		msg += "\nФраза успешно завершена!"
		msg += fmt.Sprintf("\nКоличество ошибок: %d", countErrors)
		err := a.Bot.Send(cmd.GetUserId(), msg)

		if err != nil {
			return err
		}

		session.AddData(ErrorsToHelp, 0)
		session.AddData(CountErrors, 0)
		a.ActionSession.UpdateSession(session)

		return a.newPhrase(puzzle, mode)
	}

	if mode == TransMode {
		msg += fmt.Sprintf("\nПеревод: %s", puzzleRes.Phrase.RussianText)
	}

	errorsToHelp := session.GetIntData(ErrorsToHelp)

	if puzzleRes.Result.IsCorrectAnswer {
		errorsToHelp = 0
	} else {
		msg += "\nНекорректное слово!"
		errorsToHelp += 1
		countErrors += 1
	}

	if errorsToHelp >= maxCountErrors {
		msg += fmt.Sprintf("\nСледующее слово: %s", puzzleRes.Result.NextAnswer)
	}

	session.AddData(ErrorsToHelp, errorsToHelp)
	session.AddData(CountErrors, countErrors)
	a.ActionSession.UpdateSession(session)

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
