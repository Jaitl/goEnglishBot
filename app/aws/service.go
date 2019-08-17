package aws

import (
	"errors"
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
	"time"
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
	pollySes := polly.New(sessWithRegion)
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
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(text), VoiceId: aws.String("Matthew")}

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

func (s *Session) S3UploadVoice(path, filename string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	result, err := s.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.commonSettings.S3BucketName),
		Key:    aws.String(filepath.Join(s.commonSettings.S3VoicePath, filename)),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (s *Session) S3DownloadFile(url, savePath string) error {
	file, err := os.Create(savePath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, key := ParseUrl(url)
	_, err = s.s3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &s.commonSettings.S3BucketName,
			Key:    &key,
		})

	return err
}

func (s *Session) S3DeleteFile(url string) error {
	_, key := ParseUrl(url)
	input := &s3.DeleteObjectInput{Bucket: aws.String(s.commonSettings.S3BucketName), Key: &key}
	_, err := s.svc.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) TranscribeVoice(s3Path, fileName string) (string, error) {
	media := &transcribeservice.Media{
		MediaFileUri: &s3Path,
	}

	code := "en-US"
	format := "mp3"

	input := &transcribeservice.StartTranscriptionJobInput{
		LanguageCode:         &code,
		Media:                media,
		MediaFormat:          &format,
		OutputBucketName:     &s.commonSettings.S3BucketName,
		TranscriptionJobName: &fileName,
	}

	_, err := s.transcribe.StartTranscriptionJob(input)

	if err != nil {
		return "", err
	}

	inProgress := true
	var jobRes *transcribeservice.GetTranscriptionJobOutput

	pause, _ := time.ParseDuration("10s")

	for inProgress {
		input := &transcribeservice.GetTranscriptionJobInput{TranscriptionJobName: &fileName}
		jobRes, err = s.transcribe.GetTranscriptionJob(input)

		if err != nil {
			return "", err
		}

		switch *jobRes.TranscriptionJob.TranscriptionJobStatus {
		case transcribeservice.TranscriptionJobStatusInProgress:
			inProgress = true
			time.Sleep(pause)
		case transcribeservice.TranscriptionJobStatusCompleted:
			inProgress = false
		case transcribeservice.TranscriptionJobStatusFailed:
			inProgress = false
		}
	}

	if *jobRes.TranscriptionJob.TranscriptionJobStatus == transcribeservice.TranscriptionJobStatusCompleted {
		return *jobRes.TranscriptionJob.Transcript.TranscriptFileUri, nil
	} else {
		return "", errors.New(*jobRes.TranscriptionJob.FailureReason)
	}
}
