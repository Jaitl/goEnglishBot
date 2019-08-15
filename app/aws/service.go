package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"path/filepath"
)

type Session struct {
	session        *session.Session
	translate      *translate.Translate
	polly          *polly.Polly
	commonSettings *settings.CommonSettings
}

func New(accessKey, secretKey string, commonSettings *settings.CommonSettings) (*Session, error) {
	config := aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}
	sess, err := session.NewSession(&config)

	if err != nil {
		return nil, err
	}

	trans := translate.New(sess)
	pollySes := polly.New(sess)

	return &Session{session: sess, translate: trans, polly: pollySes, commonSettings: commonSettings}, nil
}

func (s *Session) Translate(text string) (string, error) {
	input := translate.TextInput{
		SourceLanguageCode: aws.String("en"),
		TargetLanguageCode: aws.String("ru"),
		Text:               aws.String(text),
	}

	req, resp := s.translate.TextRequest(&input)

	err := req.Send()

	if err != nil {
		return "", err
	}

	return *resp.TranslatedText, nil
}

func (s *Session) Speech(text string) (string, error) {
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(text), VoiceId: aws.String("Matthew")}

	output, err := s.polly.SynthesizeSpeech(input)

	if err != nil {
		return "", err
	}

	defer output.AudioStream.Close()

	mp3FileName := uuid.Must(uuid.NewV4()).String() + ".mp3"
	mp3FilePath := filepath.Join(s.commonSettings.TmpFolder, mp3FileName)

	mp3OutFile, err := os.Create(mp3FilePath)

	if err != nil {
		return "", err
	}

	defer mp3OutFile.Close()

	_, err = io.Copy(mp3OutFile, output.AudioStream)

	if err != nil {
		return "", err
	}

	return mp3FilePath, nil
}
