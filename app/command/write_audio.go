package command

type WriteAudioCommand struct {
	UserId    int
	IncNumber int
}

func (c *WriteAudioCommand) GetUserId() int {
	return c.UserId
}

func (c *WriteAudioCommand) GetType() Type {
	return WriteAudio
}

