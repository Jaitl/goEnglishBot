package command

type PuzzleAudioCommand struct {
	UserId    int
	IncNumber int
}

func (c *PuzzleAudioCommand) GetUserId() int {
	return c.UserId
}

func (c *PuzzleAudioCommand) GetType() Type {
	return PuzzleAudio
}
