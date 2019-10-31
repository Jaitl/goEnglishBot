package command

type AddCommand struct {
	UserId int
	Text   string
}

func (c *AddCommand) GetUserId() int {
	return c.UserId
}

func (c *AddCommand) GetType() Type {
	return Add
}
