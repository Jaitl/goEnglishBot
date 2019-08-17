package aws

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	result, err := TranscribeFileParser(filepath.Join("tests", "aws_transcript.json"))

	assert.Equal(t, err, nil)
	assert.Equal(t, result, "What's up, man? What they")
}
