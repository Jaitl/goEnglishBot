package telegram

import (
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"log"
	"os"
)

type AudioService struct {
	Bot         *Telegram
	PhraseModel *phrase.Model
	AwsSession  *aws.Session
}

func NewAudioService(bot *Telegram, phModel *phrase.Model, awsSess *aws.Session) *AudioService {
	return &AudioService{
		Bot:         bot,
		PhraseModel: phModel,
		AwsSession:  awsSess,
	}
}

func (a *AudioService) SendAudio(phrs *phrase.Phrase) error {
	if phrs.AudioId != "" {
		log.Println("[DEBUG][AudioService] Send audio from cache")
		return a.Bot.SendAudioId(phrs.UserId, phrs.AudioId)
	}

	fileName, err := phrs.Title()

	if err != nil {
		return err
	}

	pathToAudioFile, err := a.AwsSession.Speech(phrs.EnglishText, fileName)

	if err != nil {
		return err
	}

	defer os.Remove(pathToAudioFile)

	log.Println("[DEBUG][AudioService] Upload new audio file")

	audioId, err := a.Bot.SendAudio(phrs.UserId, pathToAudioFile)

	if err != nil {
		return err
	}

	if audioId != "" {
		err := a.PhraseModel.UpdateAudioId(phrs.Id, audioId)

		if err != nil {
			return err
		}
	}

	return nil
}
