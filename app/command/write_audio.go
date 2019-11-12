package command

type WriteAudioCommand struct {
	UserId    int
	From   *int
	To     *int
}

func (c *WriteAudioCommand) GetUserId() int {
	return c.UserId
}

func (c *WriteAudioCommand) GetType() Type {
	return WriteAudio
}

