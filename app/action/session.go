package action

type SessionKey string

type Session struct {
	UserId     int                   `bson:"_id"`
	ActionType Type                  `bson:"action"`
	Stage      Stage                 `bson:"stage"`
	Data       map[SessionKey]string `bson:"data"`
}

type SessionModel struct {
	session map[int]*Session
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

func NewInMemorySessionModel() *SessionModel {
	return &SessionModel{session: make(map[int]*Session)}
}

func (m *SessionModel) FindSession(userId int) *Session {
	ses, ok := m.session[userId]

	if !ok {
		return nil
	}

	return ses
}

func (m *SessionModel) ClearSession(userId int) {
	delete(m.session, userId)
}

func (m *SessionModel) UpdateSession(session *Session) {
	m.session[session.UserId] = session
}
