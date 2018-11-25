package phrase

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Phrase struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	UserId      int           `bson:"userId"`
	IncNumber   int           `bson:"incNumber"`
	EnglishText string        `bson:"englishText"`
	RussianText string        `bson:"russianText"`
	IsMemorized bool          `bson:"isMemorized"`
}

type Model struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewModel(session *mgo.Session, db string) *Model {
	c := session.DB(db).C("phrase")

	return &Model{session: session, collection: c}
}

func (model *Model) CreatePhrase(userId, incNumber int, textEnglish, textRussian string) error {
	err := model.collection.Insert(Phrase{
		UserId:      userId,
		IncNumber:   incNumber,
		EnglishText: textEnglish,
		RussianText: textRussian,
		IsMemorized: false,
	})
	return err
}

func (model *Model) AllPhrases(userId int) ([]Phrase, error) {
	var phrases []Phrase

	err := model.collection.Find(bson.M{"isMemorized": false, "userId": userId}).All(&phrases)

	if err != nil {
		return nil, err
	}

	return phrases, nil
}

func (model *Model) NextIncNumber(userId int) (int, error) {
	var phrase Phrase

	err := model.collection.Find(bson.M{"isMemorized": false, "userId": userId}).Sort("-incNumber").One(&phrase)

	if err == mgo.ErrNotFound {
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return phrase.IncNumber + 1, nil
}

func (model *Model) FindPhraseByIncNumber(userId, incNumber int) (*Phrase, error) {
	var phrase Phrase

	err := model.collection.Find(bson.M{"incNumber": incNumber, "userId": userId}).One(&phrase)

	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &phrase, nil
}
