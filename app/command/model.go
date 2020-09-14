package command

type Type string

const (
	// phrase
	ListPhrases  Type = "listPhrases"
	RemovePhrase Type = "removePhrase"
	// category
	AddCategory    Type = "addCategory"
	ListCategories Type = "listCategories"
	SetCategory    Type = "setCategories"
	RemoveCategory Type = "removeCategory"
	// common
	Text             Type = "text"
	KeyboardCallback Type = "keyboardCallback"
	ReceivedVoice    Type = "receivedVoice"
	Number           Type = "number"
	// training
	PuzzleAudio Type = "puzzleAudio"
	PuzzleTrans Type = "puzzleTrans"
	WriteAudio  Type = "writeAudio"
	WriteTrans  Type = "writeTrans"
	Speech      Type = "speech"
	Card        Type = "card"
	Skip        Type = "skip"
	// system
	Me Type = "me"
)

type Command interface {
	GetUserId() int
	GetType() Type
}
