package command

type KeyboardCallbackCommand struct {
	UserId int
	Data   string
}

func (c *KeyboardCallbackCommand) GetUserId() int {
	return c.UserId
}

func (c *KeyboardCallbackCommand) GetType() Type {
	return KeyboardCallback
}
