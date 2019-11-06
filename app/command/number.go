package command

type NumberCommand struct {
	UserId    int
	IncNumber int
}

func (c *NumberCommand) GetUserId() int {
	return c.UserId
}

func (c *NumberCommand) GetType() Type {
	return Number
}
