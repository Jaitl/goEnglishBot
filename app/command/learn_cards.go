package command

type LearnCardsCommand struct {
	UserId int
	From   *int
	To     *int
}

func (c *LearnCardsCommand) GetUserId() int {
	return c.UserId
}

func (c *LearnCardsCommand) GetType() Type {
	return LearnCards
}
