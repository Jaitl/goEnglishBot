package voice

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/telegram"
)

type Action struct {
	Speech        *telegram.SpeechService
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
}

const (
	Start action.Stage = "start" // Получает id фразы
)

func (a *Action) GetType() action.Type {
	return action.Voice
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.ReceivedVoice: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		if _, ok := cmd.(*command.ReceivedVoiceCommand); ok {
			return a.voiceStage(cmd)
		}
	}

	return fmt.Errorf("stage %s not found in AddAction", stage)
}

func (a *Action) voiceStage(cmd command.Command) error {
	voiceCmd := cmd.(*command.ReceivedVoiceCommand)

	a.ActionSession.ClearSession(cmd.GetUserId())

	voiceTest, err := a.Speech.TranscribeVoice(voiceCmd.FileID)

	if err != nil {
		return err
	}

	err = a.Bot.Send(voiceCmd.GetUserId(), "Вы сказали: "+*voiceTest)

	return err
}
