package goOggCnv

import (
	"os/exec"
)

type GoOggCnv struct {
	ffmpegPath string
}

func New(ffmpegPath string) *GoOggCnv {
	return &GoOggCnv {
		ffmpegPath: ffmpegPath,
	}
}

func NewD() *GoOggCnv {
	return &GoOggCnv {
		ffmpegPath: "ffmpeg",
	}
}

func (c *GoOggCnv) Mp3ToOgg(input, output string) error {
	_, err := exec.Command(c.ffmpegPath , "-i", input, "-c:a", "libopus", output).Output()

	return err
}

func (c *GoOggCnv) WavToOgg(input, output string) error {
	_, err := exec.Command(c.ffmpegPath , "-i", input, "-acodec", "libopus", output).Output()

	return err
}
