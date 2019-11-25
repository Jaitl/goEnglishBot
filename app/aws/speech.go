package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

const (
	fileFieldName string = "file"
	rateFieldName string = "rate"
)

type fileRecognizeResult struct {
	Result string `json:"result"`
}

type SpeechClient struct {
	speechUrl string
}

func NewSpeechClient(url string) *SpeechClient {
	return &SpeechClient{speechUrl: url}
}

func (s *SpeechClient) RecognizeFile(path string, rate int) (*string, error) {
	byteJson, err := s.sendTranscribeRequest(path, rate)

	if err != nil {
		return nil, err
	}

	return s.parseJsonBytes(byteJson)
}

func (s *SpeechClient) sendTranscribeRequest(path string, rate int) ([]byte, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	r, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	// file field
	fw, err := w.CreateFormFile(fileFieldName, r.Name())

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}

	// rate field
	err = w.WriteField(rateFieldName, strconv.Itoa(rate))
	if err != nil {
		return nil, err
	}

	w.Close()

	url := s.speechUrl + "/transcribe"

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	byteRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return byteRes, nil
}

func (s *SpeechClient) parseJsonBytes(jsonArr []byte) (*string, error) {
	var data fileRecognizeResult

	err := json.Unmarshal(jsonArr, &data)

	if err != nil {
		return nil, err
	}

	return &data.Result, nil
}
