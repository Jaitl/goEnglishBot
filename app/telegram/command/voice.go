package command

type VoiceCommand struct {
	UserId    int
	IncNumber int
}

func (c *VoiceCommand) GetUserId() int {
	return c.UserId
}

func (c *VoiceCommand) GetType() Type {
	return Voice
}
