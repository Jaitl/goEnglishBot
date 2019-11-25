package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/jaitl/goEnglishBot/app/settings"
	"io"
	"os"
	"path/filepath"
)

type Session struct {
	session        *session.Session
	translate      *translate.Translate
	polly          *polly.Polly
	s3Uploader     *s3manager.Uploader
	s3Downloader   *s3manager.Downloader
	svc            *s3.S3
	transcribe     *transcribeservice.TranscribeService
	commonSettings *settings.CommonSettings
}

func New(accessKey, secretKey string, commonSettings *settings.CommonSettings) (*Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	sessWithRegion, err := session.NewSession(&aws.Config{
		Region:      aws.String(commonSettings.AwsRegion),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	trans := translate.New(sess)
	pollySes := polly.New(sess)
	transcribe := transcribeservice.New(sessWithRegion)
	s3Uploader := s3manager.NewUploader(sessWithRegion)
	s3Downloader := s3manager.NewDownloader(sessWithRegion)
	svc := s3.New(sessWithRegion)

	awsSess := Session{
		session:        sessWithRegion,
		translate:      trans,
		polly:          pollySes,
		s3Uploader:     s3Uploader,
		s3Downloader:   s3Downloader,
		svc:            svc,
		transcribe:     transcribe,
		commonSettings: commonSettings,
	}

	return &awsSess, nil
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

func (s *Session) Speech(text, name string) (string, error) {
	input := &polly.SynthesizeSpeechInput{
		Engine:       aws.String(polly.EngineNeural),
		OutputFormat: aws.String(polly.OutputFormatMp3),
		Text:         aws.String(text),
		VoiceId:      aws.String(polly.VoiceIdMatthew),
	}

	output, err := s.polly.SynthesizeSpeech(input)

	if err != nil {
		return "", err
	}

	defer output.AudioStream.Close()

	mp3FileName := name + ".mp3"
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
