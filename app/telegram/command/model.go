package command

type Type string

const (
	Add              Type = "add"
	List             Type = "list"
	Text             Type = "text"
	KeyboardCallback Type = "keyboardCallback"
	ReceivedVoice    Type = "receivedVoice"
	Audio            Type = "audio"
	Voice            Type = "voice"
)

type Command interface {
	GetUserId() int
	GetType() Type
}

type AddCommand struct {
	UserId int
	Text   string
}

type ListCommand struct {
	UserId int
}

type TextCommand struct {
	UserId int
	Text   string
}

type KeyboardCallbackCommand struct {
	UserId int
	Data   string
}

type ReceivedVoiceCommand struct {
	UserId int
	FileID string
}

type AudioCommand struct {
	UserId    int
	IncNumber int
}

type VoiceCommand struct {
	UserId    int
	IncNumber int
}

func (c *AddCommand) GetUserId() int {
	return c.UserId
}

func (c *AddCommand) GetType() Type {
	return Add
}

func (c *ListCommand) GetUserId() int {
	return c.UserId
}

func (c *ListCommand) GetType() Type {
	return List
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

func (c *ReceivedVoiceCommand) GetUserId() int {
	return c.UserId
}

func (c *ReceivedVoiceCommand) GetType() Type {
	return ReceivedVoice
}

func (c *AudioCommand) GetUserId() int {
	return c.UserId
}

func (c *AudioCommand) GetType() Type {
	return Audio
}

func (c *VoiceCommand) GetUserId() int {
	return c.UserId
}

func (c *VoiceCommand) GetType() Type {
	return Voice
}
