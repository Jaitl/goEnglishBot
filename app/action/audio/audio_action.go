package audio

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
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
	return action.Audio
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.Audio: true}
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
	audCmd := cmd.(*command.AudioCommand)

	phrs, err := a.PhraseModel.FindPhraseByIncNumber(audCmd.GetUserId(), audCmd.IncNumber)

	if err != nil {
		return err
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

	err = a.Bot.Send(audCmd.UserId, phrs.EnglishText)
	if err != nil {
		return err
	}

	err = a.Bot.SendAudio(audCmd.UserId, pathToAudioFile)

	return err
}
