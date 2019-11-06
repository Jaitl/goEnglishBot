package card

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"log"
	"os"
)

type Action struct {
	Bot         *telegram.Telegram
	PhraseModel *phrase.Model
	AwsSession  *aws.Session
}

const (
	Start action.Stage = "start"
)

func (a *Action) GetType() action.Type {
	return action.Card
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Number: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	}

	return fmt.Errorf("stage %s not found in AudioAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	audCmd := cmd.(*command.NumberCommand)

	phrs, err := a.PhraseModel.FindPhraseByIncNumber(audCmd.GetUserId(), audCmd.IncNumber)

	if err != nil {
		return err
	}

	err = a.Bot.SendMarkdown(audCmd.UserId, phrase.ToMarkdown(phrs))
	if err != nil {
		return err
	}

	if phrs.AudioId != "" {
		log.Println("[DEBUG] Send audio from cache")
		return a.Bot.SendAudioId(audCmd.UserId, phrs.AudioId)
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

	log.Println("[DEBUG] Upload new audio file")

	audioId, err := a.Bot.SendAudio(audCmd.UserId, pathToAudioFile)

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
