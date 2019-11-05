package command

type RemoveCommand struct {
	UserId    int
	IncNumber int
}

func (c *RemoveCommand) GetUserId() int {
	return c.UserId
}

func (c *RemoveCommand) GetType() Type {
	return Remove
}
