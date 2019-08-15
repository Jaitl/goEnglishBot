package main

import (
	"github.com/globalsign/mgo"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/action/add"
	"github.com/jaitl/goEnglishBot/app/action/audio"
	"github.com/jaitl/goEnglishBot/app/action/list"
	"github.com/jaitl/goEnglishBot/app/action/voice"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jessevdk/go-flags"
	"log"
)

var opts struct {
	TelegramToken   string `long:"token" env:"TOKEN" required:"true"`
	AWSKey          string `long:"aws-key" env:"AWS_KEY" required:"true"`
	AWSSecret       string `long:"aws-secret" env:"AWS_SECRET" required:"true"`
	PathToTmpFolder string `long:"tmp-folder" env:"TMP_FOLDER" required:"true"`
}

func main() {
	log.Println("[INFO] start goEnglishBot")

	if _, err := flags.Parse(&opts); err != nil {
		log.Panic(err)
	}

	commonSettings := &settings.CommonSettings{TmpFolder: opts.PathToTmpFolder}

	mongoSession, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Panic(err)
	}

	phraseModel := phrase.NewModel(mongoSession, "goEnglishBot")
	actionSession := action.NewSessionMongoModel(mongoSession, "goEnglishBot")

	awsSession, err := aws.New(opts.AWSKey, opts.AWSSecret, commonSettings)

	if err != nil {
		log.Panic(err)
	}

	telegramBot, err := telegram.New(opts.TelegramToken)

	if err != nil {
		log.Panic(err)
	}

	actions := []action.Action{
		&add.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel},
		&list.Action{Bot: telegramBot, PhraseModel: phraseModel},
		&audio.Action{Bot: telegramBot, PhraseModel: phraseModel, AwsSession: awsSession},
		&voice.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel},
	}

	actionExecutor := action.NewExecutor(actionSession, actions)

	telegramBot.Start(actionExecutor)
}
