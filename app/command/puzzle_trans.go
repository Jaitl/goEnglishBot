package command

type PuzzleTransCommand struct {
	UserId int
	From   *int
	To     *int
}

func (c *PuzzleTransCommand) GetUserId() int {
	return c.UserId
}

func (c *PuzzleTransCommand) GetType() Type {
	return PuzzleTrans
}
