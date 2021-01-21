package main

import (
	"log"

	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/action/category_crud"
	"github.com/jaitl/goEnglishBot/app/action/learn_cards"
	"github.com/jaitl/goEnglishBot/app/action/me"
	"github.com/jaitl/goEnglishBot/app/action/phrase_add"
	"github.com/jaitl/goEnglishBot/app/action/phrase_card"
	"github.com/jaitl/goEnglishBot/app/action/phrase_list"
	"github.com/jaitl/goEnglishBot/app/action/phrase_remove"
	"github.com/jaitl/goEnglishBot/app/action/puzzle"
	"github.com/jaitl/goEnglishBot/app/action/speech"
	"github.com/jaitl/goEnglishBot/app/action/voice"
	"github.com/jaitl/goEnglishBot/app/action/write"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/category"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jaitl/goEnglishBot/app/utils"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	TelegramToken    string `long:"token" env:"TOKEN" required:"true"`
	UserId           int    `long:"user-id" env:"USER_ID" required:"true"`
	MongoDbUrl       string `long:"mongo-db-url" env:"MONGO_DB_URL" required:"true"`
	AWSKey           string `long:"aws-key" env:"AWS_KEY" required:"true"`
	AWSSecret        string `long:"aws-secret" env:"AWS_SECRET" required:"true"`
	AWSRegion        string `long:"aws-region" env:"AWS_REGION" required:"true"`
	PathToTmpFolder  string `long:"tmp-folder" env:"TMP_FOLDER" required:"true"`
	SpeechServiceUrl string `long:"speech-service-url" env:"SPEECH_SERVICE_URL" required:"true"`
}

func main() {
	log.Println("[INFO] start goEnglishBot")

	if _, err := flags.Parse(&opts); err != nil {
		log.Panic(err)
	}

	commonSettings := &settings.CommonSettings{
		TmpFolder: opts.PathToTmpFolder,
		AwsRegion: opts.AWSRegion,
	}

	client, err := utils.ConnectMongo(opts.MongoDbUrl)

	if err != nil {
		log.Panic(err)
	}

	categoryModel, err := category.NewModel(client, "goEnglishBot")

	if err != nil {
		log.Panic(err)
	}

	actionSession := action.NewInMemorySessionModel()

	awsSession, err := aws.New(opts.AWSKey, opts.AWSSecret, commonSettings)

	if err != nil {
		log.Panic(err)
	}

	telegramBot, err := telegram.New(opts.TelegramToken, opts.UserId)

	if err != nil {
		log.Panic(err)
	}

	audioService := telegram.NewAudioService(telegramBot, categoryModel, awsSession)

	speechService := telegram.NewSpeechService(telegramBot, awsSession, commonSettings)

	actions := []action.Action{
		&phrase_add.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, CategoryModel: categoryModel},
		&phrase_list.Action{Bot: telegramBot, CategoryModel: categoryModel},
		&phrase_remove.Action{Bot: telegramBot, CategoryModel: categoryModel},
		&phrase_card.Action{Bot: telegramBot, CategoryModel: categoryModel, AwsSession: awsSession, Audio: audioService},
		&category_crud.Action{Bot: telegramBot, CategoryModel: categoryModel},
		&voice.Action{Speech: speechService, ActionSession: actionSession, Bot: telegramBot},
		&me.Action{Bot: telegramBot},
		&puzzle.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, CategoryModel: categoryModel, Audio: audioService},
		&write.Action{ActionSession: actionSession, Bot: telegramBot, CategoryModel: categoryModel, Audio: audioService},
		&speech.Action{ActionSession: actionSession, Bot: telegramBot, CategoryModel: categoryModel, Speech: speechService, Audio: audioService},
		&learn_cards.Action{Bot: telegramBot, CategoryModel: categoryModel, AwsSession: awsSession, Audio: audioService, ActionSession: actionSession},
	}

	actionExecutor := action.NewExecutor(actionSession, actions)

	telegramBot.Start(actionExecutor)
}
