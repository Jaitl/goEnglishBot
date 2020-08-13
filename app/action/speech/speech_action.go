package speech

import (
	"errors"
	"fmt"
	"time"

	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/exercises"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/utils"
)

const (
	Start     action.Stage = "start"
	WaitVoice action.Stage = "waitVoice"
)

const (
	Session     action.SessionKey = "speechSession"
	CountErrors action.SessionKey = "countErrors"
	StartTime   action.SessionKey = "startTime"
)

type Action struct {
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	CategoryModel *category.Model
	Speech        *telegram.SpeechService
	Audio         *telegram.AudioService
}

func (a *Action) GetType() action.Type {
	return action.Speech
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Speech: true}
	case WaitVoice:
		return map[command.Type]bool{command.ReceivedVoice: true, command.Skip: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitVoice:
		if cm, ok := cmd.(*command.SkipCommand); ok {
			return a.skipPhrase(cm, session)
		} else {
			return a.waitVoice(cmd, session)
		}
	}

	return fmt.Errorf("stage %s not found in SpeechAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	speechCmd := cmd.(*command.SpeechCommand)

	phrs, err := a.CategoryModel.SmartFindByRange(cmd.GetUserId(), speechCmd.From, speechCmd.To)

	if err != nil {
		return err
	}

	if len(phrs) == 0 {
		return errors.New("don't correct range")
	}

	speech := exercises.NewComposite(phrs, exercises.SpeechMode, true)

	err = a.newSpeech(speech, false)

	ses := action.CreateSession(cmd.GetUserId(), action.Speech, WaitVoice)
	ses.AddData(Session, speech)
	ses.AddData(CountErrors, 0)
	ses.AddData(StartTime, time.Now())
	a.ActionSession.UpdateSession(ses)

	return err
}

func (a *Action) waitVoice(cmd command.Command, session *action.Session) error {
	voice, ok := cmd.(*command.ReceivedVoiceCommand)

	if !ok {
		return errors.New("command does not belong to WaitPushButton stage in PuzzleAction")
	}

	speech := session.Data[Session].(*exercises.Composite)
	countErrors := session.GetIntData(CountErrors)

	answerText, err := a.Speech.TranscribeVoice(voice.FileID)

	if err != nil {
		return err
	}

	text := exercises.ClearText(*answerText)

	speechRes := speech.HandleAnswer(exercises.ClearText(*answerText))

	msg := fmt.Sprintf("Фраза №%d из %d", speechRes.Pos+1, speechRes.CountPhrases)
	msg += fmt.Sprintf("\nФраза: %s", exercises.ClearText(speechRes.Phrase.EnglishText))
	msg += fmt.Sprintf("\nВы сказали: %s", text)
	msg += fmt.Sprintf("\nСовпадение: %d%%", int(speechRes.Result.MatchScore*100))

	if speechRes.Result.IsFinish && speechRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		msg += fmt.Sprintf("\nПеревод: %s", speechRes.Phrase.RussianText)
		msg += "\nФраза успешно завершена!"
		msg += fmt.Sprintf("\nКоличество ошибок: %d", countErrors)

		startTime := session.Data[StartTime].(time.Time)
		elapsed := time.Since(startTime)

		msg += fmt.Sprintf("\nУпражнение успешно завершено за: %s!", utils.DurationPretty(elapsed))

		return a.Bot.Send(cmd.GetUserId(), msg)
	}

	if speechRes.Result.IsFinish && !speechRes.IsFinish {
		msg += fmt.Sprintf("\nПеревод: %s", speechRes.Phrase.RussianText)
		msg += "\nФраза успешно завершена!"
		msg += fmt.Sprintf("\nКоличество ошибок: %d", countErrors)
		err := a.Bot.Send(cmd.GetUserId(), msg)

		if err != nil {
			return err
		}

		session.AddData(CountErrors, 0)
		a.ActionSession.UpdateSession(session)

		return a.newSpeech(speech, false)
	}

	if !speechRes.Result.IsCorrectAnswer {
		countErrors += 1
		msg += fmt.Sprintf("\nПроизношение некорректно! Отправьте голосовое сообщение еще раз")
	}

	session.AddData(CountErrors, countErrors)
	a.ActionSession.UpdateSession(session)

	return a.Bot.Send(cmd.GetUserId(), msg)
}

func (a *Action) skipPhrase(cmd *command.SkipCommand, session *action.Session) error {
	session.AddData(CountErrors, 0)
	a.ActionSession.UpdateSession(session)
	speech := session.Data[Session].(*exercises.Composite)
	speechRes := speech.Skip()

	if speechRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		startTime := session.Data[StartTime].(time.Time)
		elapsed := time.Since(startTime)

		msg := fmt.Sprintf("Упражнение успешно завершено за: %s!", utils.DurationPretty(elapsed))

		return a.Bot.Send(cmd.GetUserId(), msg)
	}

	return a.newSpeech(speech, true)
}

func (a *Action) newSpeech(speech *exercises.Composite, skip bool) error {
	speechRes := speech.Next()

	msg := fmt.Sprintf("Фраза №%d из %d", speechRes.Pos+1, speechRes.CountPhrases)

	msg += "\nОтправьте голосовое сообщение с произношением фразы:"
	msg += "\n" + exercises.ClearText(speechRes.Phrase.EnglishText)

	err := a.Bot.SendMarkdown(speechRes.Phrase.UserId, msg)
	if err != nil {
		return err
	}

	return a.Audio.SendAudio(speechRes.Phrase)
}
