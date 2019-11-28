package command

// list
type ListPhrasesCommand struct {
	UserId    int
	IncNumber *int
}

func (c *ListPhrasesCommand) GetUserId() int {
	return c.UserId
}

func (c *ListPhrasesCommand) GetType() Type {
	return ListPhrases
}

// remove
type RemovePhraseCommand struct {
	UserId    int
	IncNumber int
}

func (c *RemovePhraseCommand) GetUserId() int {
	return c.UserId
}

func (c *RemovePhraseCommand) GetType() Type {
	return RemovePhrase
}
