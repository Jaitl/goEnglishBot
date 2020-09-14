package command

type CardCommand struct {
	UserId int
	From   *int
	To     *int
}

func (c *CardCommand) GetUserId() int {
	return c.UserId
}

func (c *CardCommand) GetType() Type {
	return Card
}
