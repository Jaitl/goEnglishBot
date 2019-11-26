package command

type Type string

const (
	Add              Type = "add"
	List             Type = "list"
	Text             Type = "text"
	KeyboardCallback Type = "keyboardCallback"
	ReceivedVoice    Type = "receivedVoice"
	Number           Type = "number"
	Speech           Type = "speech"
	Me               Type = "me"
	Remove           Type = "remove"
	PuzzleAudio      Type = "puzzleAudio"
	PuzzleTrans      Type = "puzzleTrans"
	WriteAudio       Type = "writeAudio"
	WriteTrans       Type = "writeTrans"
)

type Command interface {
	GetUserId() int
	GetType() Type
}
