package telegram

import (
	"log"
	"net/http"

	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/settings"
)

type SpeechService struct {
	bot            *Telegram
	awsSession     *aws.Session
	commonSettings *settings.CommonSettings
}

func NewSpeechService(bot *Telegram, awsSession *aws.Session, sett *settings.CommonSettings) *SpeechService {
	return &SpeechService{
		bot:            bot,
		awsSession:     awsSession,
		commonSettings: sett,
	}
}

func (a *SpeechService) TranscribeVoice(fileId string) (*string, error) {
	file, err := a.bot.bot.GetFileDirectURL(fileId)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(file)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("[DEBUG] [SpeechService]: Do request to recognize voice")

	return a.awsSession.Transcribe(resp.Body)
}
