package telegram

import (
	"github.com/google/uuid"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	rate int = 16000
)

type SpeechService struct {
	bot            *Telegram
	speech         *aws.SpeechClient
	commonSettings *settings.CommonSettings
}

func NewSpeechService(bot *Telegram, speech *aws.SpeechClient, sett *settings.CommonSettings) *SpeechService {
	return &SpeechService{
		bot:            bot,
		speech:         speech,
		commonSettings: sett,
	}
}

func (a *SpeechService) TranscribeVoice(fileId string) (*string, error) {
	fileUrl, err := a.bot.GetFilePath(fileId)

	if err != nil {
		return nil, err
	}

	opusFileTmpUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	opusFileTmpName := opusFileTmpUuid.String() + ".opus"
	opusFilePath := filepath.Join(a.commonSettings.TmpFolder, opusFileTmpName)

	log.Println("[DEBUG] [SpeechService]: Download file from Telegram")
	err = utils.DownloadFile(opusFilePath, fileUrl)

	if err != nil {
		return nil, err
	}

	defer os.Remove(opusFilePath)

	pcmFileTmpUuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	pcmFileTmpName := pcmFileTmpUuid.String() + ".pcm"
	pcmFilePath := filepath.Join(a.commonSettings.TmpFolder, pcmFileTmpName)

	log.Println("[DEBUG] [SpeechService]: Convert file to PCM")
	err = utils.OpusToPcm(opusFilePath, pcmFilePath, strconv.Itoa(rate))

	if err != nil {
		return nil, err
	}

	defer os.Remove(pcmFilePath)

	log.Println("[DEBUG] [SpeechService]: Do request to recognize voice")

	return a.speech.RecognizeFile(pcmFilePath, rate)
}
