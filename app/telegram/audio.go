package telegram

import (
	"log"

	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/phrase"

	tgbotgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type AudioService struct {
	Bot           *Telegram
	CategoryModel *category.Model
	AwsSession    *aws.Session
}

func NewAudioService(bot *Telegram, catModel *category.Model, awsSess *aws.Session) *AudioService {
	return &AudioService{
		Bot:           bot,
		CategoryModel: catModel,
		AwsSession:    awsSess,
	}
}

func (a *AudioService) SendAudio(phrs *phrase.Phrase) error {
	if phrs.AudioId != "" {
		log.Println("[DEBUG][AudioService] Send audio from cache")
		return a.Bot.SendAudioId(phrs.UserId, phrs.AudioId)
	}
	audioBytesArray, err := a.AwsSession.Speech(phrs.EnglishText)

	if err != nil {
		return err
	}

	log.Println("[DEBUG][AudioService] Upload new audio file")

	fileName := phrs.Title()
	fileBytes := tgbotgapi.FileBytes{Name: fileName, Bytes: audioBytesArray}

	audioId, err := a.Bot.SendAudio(phrs.UserId, fileBytes)

	if err != nil {
		return err
	}

	if audioId != "" {
		err := a.CategoryModel.UpdatePhraseAudioId(phrs.UserId, phrs.IncNumber, audioId)

		if err != nil {
			return err
		}
	}

	return nil
}
