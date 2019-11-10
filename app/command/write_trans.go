package command

type WriteTransCommand struct {
	UserId    int
	IncNumber int
}

func (c *WriteTransCommand) GetUserId() int {
	return c.UserId
}

func (c *WriteTransCommand) GetType() Type {
	return WriteTrans
}
