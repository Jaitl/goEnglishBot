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
