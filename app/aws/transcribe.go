package aws

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type transcript struct {
	Transcript string `json:"transcript"`
}

type transcripts struct {
	Transcripts []transcript `json:"transcripts"`
}

type result struct {
	Results transcripts `json:"results"`
}

func TranscribeFileParser(path string) (string, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return TranscribeJsonParser(b)
}

func TranscribeJsonParser(jsonArr []byte) (string, error) {
	var data result

	err := json.Unmarshal(jsonArr, &data)

	if err != nil {
		return "", err
	}

	var phrases []string

	results := data.Results
	for _, trans := range results.Transcripts {
		phrases = append(phrases, trans.Transcript)
	}

	return strings.Join(phrases, " "), nil
}
