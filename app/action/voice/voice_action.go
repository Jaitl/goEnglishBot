package voice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Action struct {
	Speech         *aws.SpeechClient
	ActionSession  *action.SessionModel
	Bot            *telegram.Telegram
	PhraseModel    *phrase.Model
	CommonSettings *settings.CommonSettings
}

const (
	Start action.Stage = "start" // Получает id фразы

	rate int = 16000
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

	fileUrl, err := a.Bot.GetFilePath(voiceCmd.FileID)

	if err != nil {
		return err
	}

	opusFileTmpUuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	opusFileTmpName := opusFileTmpUuid.String() + ".opus"
	opusFilePath := filepath.Join(a.CommonSettings.TmpFolder, opusFileTmpName)

	log.Println("[DEBUG] VOICE: Download file from Telegram")
	err = utils.DownloadFile(opusFilePath, fileUrl)

	if err != nil {
		return err
	}

	defer os.Remove(opusFilePath)

	pcmFileTmpUuid, err := uuid.NewRandom()

	if err != nil {
		return err
	}

	pcmFileTmpName := pcmFileTmpUuid.String() + ".pcm"
	pcmFilePath := filepath.Join(a.CommonSettings.TmpFolder, pcmFileTmpName)

	log.Println("[DEBUG] VOICE: Convert file to PCM")
	err = utils.OpusToPcm(opusFilePath, pcmFilePath, strconv.Itoa(rate))

	if err != nil {
		return err
	}

	defer os.Remove(pcmFilePath)

	log.Println("[DEBUG] VOICE: Do request to recognize voice")

	voiceTest, err := a.Speech.RecognizeFile(pcmFilePath, rate)

	if err != nil {
		return err
	}

	err = a.Bot.Send(voiceCmd.GetUserId(), "Вы сказали: "+*voiceTest)

	return err
}
