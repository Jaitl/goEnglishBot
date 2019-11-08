package command

type PuzzleTransCommand struct {
	UserId    int
	IncNumber int
}

func (c *PuzzleTransCommand) GetUserId() int {
	return c.UserId
}

func (c *PuzzleTransCommand) GetType() Type {
	return PuzzleTrans
}
