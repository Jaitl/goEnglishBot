package main

import (
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/action/add"
	"github.com/jaitl/goEnglishBot/app/action/card"
	"github.com/jaitl/goEnglishBot/app/action/list"
	"github.com/jaitl/goEnglishBot/app/action/me"
	"github.com/jaitl/goEnglishBot/app/action/puzzle"
	"github.com/jaitl/goEnglishBot/app/action/remove"
	"github.com/jaitl/goEnglishBot/app/action/speech"
	"github.com/jaitl/goEnglishBot/app/action/voice"
	"github.com/jaitl/goEnglishBot/app/action/write"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jessevdk/go-flags"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

	client, err := mongo.NewClient(options.Client().ApplyURI(opts.MongoDbUrl))

	if err != nil {
		log.Panic(err)
	}

	phraseModel, err := phrase.NewModel(client, "goEnglishBot")

	if err != nil {
		log.Panic(err)
	}

	actionSession := action.NewInMemorySessionModel()

	awsSession, err := aws.New(opts.AWSKey, opts.AWSSecret, commonSettings)

	if err != nil {
		log.Panic(err)
	}

	speechClient := aws.NewSpeechClient(opts.SpeechServiceUrl)

	telegramBot, err := telegram.New(opts.TelegramToken, opts.UserId)

	if err != nil {
		log.Panic(err)
	}

	audioService := telegram.NewAudioService(telegramBot, phraseModel, awsSession)

	speechServie := telegram.NewSpeechService(telegramBot, speechClient, commonSettings)

	actions := []action.Action{
		&add.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel},
		&list.Action{Bot: telegramBot, PhraseModel: phraseModel},
		&card.Action{Bot: telegramBot, PhraseModel: phraseModel, AwsSession: awsSession, Audio: audioService},
		&voice.Action{Speech: speechServie, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel},
		&me.Action{Bot: telegramBot},
		&remove.Action{Bot: telegramBot, PhraseModel: phraseModel},
		&puzzle.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel, Audio: audioService},
		&write.Action{ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel, Audio: audioService},
		&speech.Action{ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel, Speech: speechServie, Audio: audioService},
	}

	actionExecutor := action.NewExecutor(actionSession, actions)

	telegramBot.Start(actionExecutor)
}
