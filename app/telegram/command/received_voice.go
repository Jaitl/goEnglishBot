package command

type ReceivedVoiceCommand struct {
	UserId int
	FileID string
}

func (c *ReceivedVoiceCommand) GetUserId() int {
	return c.UserId
}

func (c *ReceivedVoiceCommand) GetType() Type {
	return ReceivedVoice
}
