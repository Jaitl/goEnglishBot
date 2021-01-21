package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	transcribe "github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/jaitl/goEnglishBot/app/settings"
)

type Session struct {
	session        *session.Session
	translate      *translate.Translate
	polly          *polly.Polly
	transcribe     *transcribe.TranscribeStreamingService
	commonSettings *settings.CommonSettings
}

func New(accessKey, secretKey string, commonSettings *settings.CommonSettings) (*Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(commonSettings.AwsRegion),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	translateSess := translate.New(sess)
	pollySess := polly.New(sess)
	transcribeSess := transcribe.New(sess)

	awsSess := Session{
		session:        sess,
		translate:      translateSess,
		polly:          pollySess,
		transcribe:     transcribeSess,
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

func (s *Session) Transcribe(path string, rate int) (*string, error) {
	resp, err := s.transcribe.StartStreamTranscription(&transcribe.StartStreamTranscriptionInput{
		LanguageCode:         aws.String(transcribe.LanguageCodeEnUs),
		MediaEncoding:        aws.String(transcribe.MediaEncodingPcm),
		MediaSampleRateHertz: aws.Int64(int64(rate)),
	})
	if err != nil {
		return nil, err
	}
	stream := resp.GetStream()
	defer stream.Close()

	audio, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer audio.Close()

	go func() {
		err := transcribe.StreamAudioFromReader(context.Background(), stream.Writer, 10*1024, audio)
		if err != nil {
			log.Printf("[ERROR] fail to start StreamAudioFromReader, err: %s", err.Error())
		}
	}()

	result := ""

	for event := range stream.Events() {
		switch e := event.(type) {
		case *transcribe.TranscriptEvent:
			transcr := ""
			// log.Printf("full: %s", e.Transcript.String())
			for _, res := range e.Transcript.Results {
				transcr = transcr + " " + aws.StringValue(res.Alternatives[0].Transcript)
			}
			result = transcr
		default:
			return nil, fmt.Errorf("unexpected event, %T", event)
		}
	}

	if err := stream.Err(); err != nil {
		return nil, err
	}

	result = strings.TrimSpace(result)

	return &result, nil
}
