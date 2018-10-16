package action

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SessionKey string

type Session struct {
	UserId     int                   `bson:"_id"`
	ActionType Type                  `bson:"action"`
	Stage      Stage                 `bson:"stage"`
	Data       map[SessionKey]string `bson:"data"`
}

type SessionModel struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func CreateSession(userId int, action Type, stage Stage) *Session {
	return &Session{
		UserId:     userId,
		ActionType: action,
		Stage:      stage,
		Data:       make(map[SessionKey]string),
	}
}

func (s *Session) AddData(key SessionKey, value string) {
	s.Data[key] = value
}

func NewSessionMongoModel(session *mgo.Session, db string) *SessionModel {
	c := session.DB(db).C("session")

	return &SessionModel{session: session, collection: c}
}

func (m *SessionModel) FindSession(userId int) (*Session, error) {
	var ses Session
	err := m.collection.Find(bson.M{"_id": userId}).One(&ses)

	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &ses, nil
}

func (m *SessionModel) ClearSession(userId int) error {
	err := m.collection.Remove(bson.M{"_id": userId})

	if err == mgo.ErrNotFound {
		return nil
	}

	return err
}

func (m *SessionModel) UpdateSession(session *Session) error {
	err := m.ClearSession(session.UserId)

	if err != nil {
		return err
	}

	err = m.collection.Insert(session)

	return err
}
