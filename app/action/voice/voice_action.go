package voice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
	"github.com/jaitl/goEnglishBot/app/utils"
	"os"
	"path/filepath"
)

type Action struct {
	AwsSession     *aws.Session
	ActionSession  *action.SessionModel
	Bot            *telegram.Telegram
	PhraseModel    *phrase.Model
	CommonSettings *settings.CommonSettings
}

const (
	Start action.Stage = "start" // Получает id фразы
	Voice action.Stage = "voice" // Получает произнесенную фразу

	phraseId   action.SessionKey = "phraseId"
	phraseText action.SessionKey = "phraseText"

	voiceMsg string = "Отправьте голосовое сообщение с произношением фразы \"%v\""
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
		return map[command.Type]bool{command.Voice: true, command.ReceivedVoice: true}
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

	err = utils.DownloadFile(opusFilePath, fileUrl)

	if err != nil {
		return err
	}

	defer os.Remove(opusFilePath)

	mp3FileTmpUuid, err := uuid.NewRandom()

	if err != nil {
		return err
	}

	mp3FileTmpName := mp3FileTmpUuid.String() + ".mp3"
	mp3FilePath := filepath.Join(a.CommonSettings.TmpFolder, mp3FileTmpName)

	err = utils.OpusToMp3(opusFilePath, mp3FilePath)

	if err != nil {
		return err
	}

	defer os.Remove(mp3FilePath)

	s3Url, err := a.AwsSession.S3UploadVoice(mp3FilePath, mp3FileTmpName)

	if err != nil {
		return err
	}

	defer func() {
		err := a.AwsSession.S3DeleteFile(s3Url)
		if err != nil {
			println(err.Error())
		}
	}()

	s3TranscribeUrl, err := a.AwsSession.TranscribeVoice(s3Url, mp3FileTmpName)

	if err != nil {
		return err
	}

	defer func() {
		err := a.AwsSession.S3DeleteFile(s3TranscribeUrl)
		if err != nil {
			println(err.Error())
		}
	}()

	transcribeFilePath := mp3FilePath + ".trans"

	err = a.AwsSession.S3DownloadFile(s3TranscribeUrl, transcribeFilePath)

	if err != nil {
		return err
	}

	defer os.Remove(transcribeFilePath)

	voiceTest, err := aws.TranscribeFileParser(transcribeFilePath)

	if err != nil {
		return err
	}

	err = a.Bot.Send(voiceCmd.GetUserId(), "Вы сказали: "+voiceTest)

	return err
}
