package phrase

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

const (
	rowPattern = "#%v \"*%v*\": _%v_"
)

type Phrase struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      int                `bson:"userId"`
	IncNumber   int                `bson:"incNumber"`
	EnglishText string             `bson:"englishText"`
	RussianText string             `bson:"russianText"`
	IsMemorized bool               `bson:"isMemorized"`
	AudioId     string             `bson:"audioId"`
}

func (p *Phrase) Title() (string, error) {
	title := "#" + strconv.Itoa(p.IncNumber)
	return title, nil
}

func (p *Phrase) ToMarkdown() string {
	return fmt.Sprintf(rowPattern, p.IncNumber, p.EnglishText, p.RussianText)
}
