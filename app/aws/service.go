package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/translate"
)

type Session struct {
	session   *session.Session
	translate *translate.Translate
	polly     *polly.Polly
}

func New(accessKey, secretKey string) (*Session, error) {
	config := aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}
	sess, err := session.NewSession(&config)

	if err != nil {
		return nil, err
	}

	trans := translate.New(sess)

	polly := polly.New(sess)

	return &Session{session: sess, translate: trans, polly: polly}, nil
}

func (s *Session) Translate(text string) (string, error) {
	input := translate.TextInput{
		SourceLanguageCode: aws.String("en"),
		TargetLanguageCode: aws.String("ru"),
		Text: aws.String(text),
	}

	req, resp := s.translate.TextRequest(&input)

	err := req.Send()

	if err != nil { return "", err }

	return *resp.TranslatedText, nil
}

func (s *Session) Speach(text string) (error) {
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("ogg_vorbis"), Text: aws.String(text), VoiceId: aws.String("Matthew")}

	output, err := s.polly.SynthesizeSpeech(input)

	output.AudioStream
}
