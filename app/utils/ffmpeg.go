package utils

import (
	"os/exec"
)

func OpusToMp3(input, output string) error {
	_, err := exec.Command("ffmpeg", "-i", input, "-acodec", "libmp3lame", output).Output()

	return err
}
