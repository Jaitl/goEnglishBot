package voice

import (
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
	Start action.Stage = "start" // Получает id фразы
	Voice action.Stage = "voice" // Получает произнесенную фразу

	phraseId action.SessionKey = "phraseId"
	phraseText action.SessionKey = "phraseText"

	voiceMsg string = "Отправьте голосовое сообщение с произношением фразы \"%v\""
)

func (a *Action) GetType() action.Type {
	return action.Voice
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetStartCommands() []command.Type {
	return []command.Type{command.Voice}
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Voice: true}
	case Voice:
		return map[command.Type]bool{command.ReceivedVoice: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case Voice:
		return a.voiceStage(cmd, session)
	}

	return fmt.Errorf("stage %s not found in AddAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	voiceCmd := cmd.(*command.VoiceCommand)

	phrs, err := a.PhraseModel.FindPhraseByIncNumber(voiceCmd.GetUserId(), voiceCmd.IncNumber)

	if err != nil {
		return err
	}

	ses := action.CreateSession(cmd.GetUserId(), action.Voice, Voice)
	ses.AddData(phraseId, string(voiceCmd.IncNumber))
	ses.AddData(phraseText, phrs.EnglishText)
	err = a.ActionSession.UpdateSession(ses)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(voiceMsg, phrs.EnglishText)
	err = a.Bot.Send(voiceCmd.GetUserId(), msg)

	return err
}

func (a *Action) voiceStage(cmd command.Command, session *action.Session) error {
	voiceCmd := cmd.(*command.ReceivedVoiceCommand)

	err := a.Bot.Send(voiceCmd.GetUserId(), "Тут будет ваша фраза")

	return err
}
