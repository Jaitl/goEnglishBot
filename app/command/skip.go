package command

type SkipCommand struct {
	UserId int
}

func (c *SkipCommand) GetUserId() int {
	return c.UserId
}

func (c *SkipCommand) GetType() Type {
	return Skip
}
