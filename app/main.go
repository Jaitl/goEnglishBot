package main

import (
	"github.com/globalsign/mgo"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/action/add"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jessevdk/go-flags"
	"log"
)

var opts struct {
	TelegramToken string `long:"token" env:"TOKEN" required:"true"`
	AWSKey        string `long:"aws-key" env:"AWS_KEY" required:"true"`
	AWSSecret     string `long:"aws-secret" env:"AWS_SECRET" required:"true"`
}

func main() {
	log.Println("[INFO] start goEnglishBot")

	if _, err := flags.Parse(&opts); err != nil {
		log.Panic(err)
	}

	mongoSession, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Panic(err)
	}

	phraseModel := phrase.New(mongoSession, "goEnglishBot")
	actionSession := action.NewSessionMongoModel(mongoSession, "goEnglishBot")

	awsSession, err := aws.New(opts.AWSKey, opts.AWSSecret)

	if err != nil {
		log.Panic(err)
	}

	telegramBot, err := telegram.New(opts.TelegramToken)

	if err != nil {
		log.Panic(err)
	}

	actions := []action.Action{
		&add.Action{awsSession, actionSession, telegramBot, phraseModel},
	}

	actionExecutor := action.NewExecutor(actionSession, actions)

	telegramBot.Start(actionExecutor)
}
