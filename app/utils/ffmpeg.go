package utils

import (
	"os/exec"
)

func OpusToMp3(input, output string) error {
	_, err := exec.Command("ffmpeg", "-i", input, "-acodec", "libmp3lame", output).Output()

	return err
}

func OpusToPcm(input, output string, rate string) error {
	_, err := exec.Command("ffmpeg", "-i", input, "-acodec", "pcm_s16le", "-f", "s16le", "-ac", "1", "-ar", rate, output).Output()

	return err
}
