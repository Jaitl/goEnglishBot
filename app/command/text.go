package command

type TextCommand struct {
	UserId int
	Text   string
}

func (c *TextCommand) GetUserId() int {
	return c.UserId
}

func (c *TextCommand) GetType() Type {
	return Text
}
