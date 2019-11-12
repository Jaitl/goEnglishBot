package command

type PuzzleAudioCommand struct {
	UserId int
	From   *int
	To     *int
}

func (c *PuzzleAudioCommand) GetUserId() int {
	return c.UserId
}

func (c *PuzzleAudioCommand) GetType() Type {
	return PuzzleAudio
}
