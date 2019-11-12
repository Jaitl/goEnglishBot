package command

type WriteTransCommand struct {
	UserId    int
	From   *int
	To     *int
}

func (c *WriteTransCommand) GetUserId() int {
	return c.UserId
}

func (c *WriteTransCommand) GetType() Type {
	return WriteTrans
}
