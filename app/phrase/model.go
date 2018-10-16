package phrase

import (
	"github.com/globalsign/mgo"
)

type Phrase struct {
	UserId    int    `bson:"userId"`
	Text      string `bson:"text"`
	Translate string `bson:"translate"`
}

type Model struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func New(session *mgo.Session, db string) *Model {
	c := session.DB(db).C("phrase")

	return &Model{session: session, collection: c}
}

func (model *Model) Create(userId int, text, translate string) error {
	err := model.collection.Insert(Phrase{UserId: userId, Text: text, Translate: translate})
	return err
}

func (model *Model) AllPhrases() ([]Phrase, error) {
	var phrases []Phrase

	err := model.collection.Find(nil).All(&phrases)

	if err != nil {
		return nil, err
	}

	return phrases, nil
}
