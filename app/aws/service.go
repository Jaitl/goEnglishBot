package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

type Session struct {
	session   *session.Session
	translate *translate.Translate
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

	return &Session{session: sess, translate: trans}, nil
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
