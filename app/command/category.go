package command

// add
type AddCategoryCommand struct {
	UserId int
	Name   string
}

func (c *AddCategoryCommand) GetUserId() int {
	return c.UserId
}

func (c *AddCategoryCommand) GetType() Type {
	return AddCategory
}

// list
type ListCategoriesCommand struct {
	UserId int
}

func (c *ListCategoriesCommand) GetUserId() int {
	return c.UserId
}

func (c *ListCategoriesCommand) GetType() Type {
	return ListCategories
}

// set
type SetCategoriesCommand struct {
	UserId    int
	IncNumber int
}

func (c *SetCategoriesCommand) GetUserId() int {
	return c.UserId
}

func (c *SetCategoriesCommand) GetType() Type {
	return SetCategory
}

// remove
type RemoveCategoryCommand struct {
	UserId    int
	IncNumber int
}

func (c *RemoveCategoryCommand) GetUserId() int {
	return c.UserId
}

func (c *RemoveCategoryCommand) GetType() Type {
	return RemoveCategory
}
