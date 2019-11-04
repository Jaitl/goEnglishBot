package phrase

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strconv"
	"strings"
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

const (
	phraseTitleSize = 40
)

func (phrase *Phrase) Title() (string, error) {
	reg := regexp.MustCompile(`[^a-zA-Z\s]+`)

	processedString := reg.ReplaceAllString(phrase.EnglishText, "")

	title := "#" + strconv.Itoa(phrase.IncNumber)

	parts := strings.Split(processedString, " ")

	for _, part := range parts {
		title = title + "-" + part
		if len(title) > phraseTitleSize {
			break
		}
	}

	return title, nil
}
