package command

type AudioCommand struct {
	UserId    int
	IncNumber int
}

func (c *AudioCommand) GetUserId() int {
	return c.UserId
}

func (c *AudioCommand) GetType() Type {
	return Audio
}
