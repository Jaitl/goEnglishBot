package command

type Type string

const (
	Add              Type = "add"
	List             Type = "list"
	Text             Type = "text"
	KeyboardCallback Type = "keyboardCallback"
	ReceivedVoice    Type = "receivedVoice"
	Number            Type = "number"
	Voice            Type = "voice"
	Me               Type = "me"
	Remove           Type = "remove"
)

type Command interface {
	GetUserId() int
	GetType() Type
}
