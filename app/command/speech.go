package command

type SpeechCommand struct {
	UserId    int
	From   *int
	To     *int
}

func (c *SpeechCommand) GetUserId() int {
	return c.UserId
}

func (c *SpeechCommand) GetType() Type {
	return Speech
}
