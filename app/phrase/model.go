package phrase

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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

func (phrase *Phrase) Title() (string, error) {
	title := "#" + strconv.Itoa(phrase.IncNumber)
	return title, nil
}
