package command

type Type string

const (
	Add              Type = "add"
	Text             Type = "text"
	KeyboardCallback Type = "keyboardCallback"
)

type Command interface {
	GetUserId() int
	GetType() Type
}

type AddCommand struct {
	UserId int
	Text   string
}

type TextCommand struct {
	UserId int
	Text   string
}

type KeyboardCallbackCommand struct {
	UserId int
	Data   string
}

func (c *AddCommand) GetUserId() int {
	return c.UserId
}

func (c *AddCommand) GetType() Type {
	return Add
}

func (c *TextCommand) GetUserId() int {
	return c.UserId
}

func (c *TextCommand) GetType() Type {
	return Text
}

func (c *KeyboardCallbackCommand) GetUserId() int {
	return c.UserId
}

func (c *KeyboardCallbackCommand) GetType() Type {
	return KeyboardCallback
}
