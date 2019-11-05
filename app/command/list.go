package command

type ListCommand struct {
	UserId int
}

func (c *ListCommand) GetUserId() int {
	return c.UserId
}

func (c *ListCommand) GetType() Type {
	return List
}
