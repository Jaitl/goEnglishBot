package command

type MeCommand struct {
	UserId int
}

func (c *MeCommand) GetUserId() int {
	return c.UserId
}

func (c *MeCommand) GetType() Type {
	return Me
}
