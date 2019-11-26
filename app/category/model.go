package category

import (
	"fmt"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	activeRowPattern    = "#%d *%s* _(size: %d, date: %s)_"
	notActiveRowPattern = "#%d %s _(size: %d, date: %s)_"
)

type Category struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Active     bool               `bson:"active"`
	UserId     int                `bson:"userId"`
	IncNumber  int                `bson:"incNumber"`
	Name       string             `bson:"name"`
	CreateDate int64              `bson:"createDate"`
	Phrases    []phrase.Phrase    `bson:"phrases"`
}

func (c *Category) ToMarkdown() string {
	createDate := time.Unix(c.CreateDate, 0)
	dateStr := createDate.Format("02.01.2006")
	if c.Active {
		return fmt.Sprintf(activeRowPattern, c.IncNumber, c.Name, len(c.Phrases), dateStr)
	} else {
		return fmt.Sprintf(notActiveRowPattern, c.IncNumber, c.Name, len(c.Phrases), dateStr)
	}
}
