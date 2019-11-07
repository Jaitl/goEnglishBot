package action

import "sync"

type SessionKey string

type Session struct {
	UserId     int
	ActionType Type
	Stage      Stage
	Data       map[SessionKey]interface{}
}

type SessionModel struct {
	session map[int]*Session
	mux     sync.Mutex
}

func CreateSession(userId int, action Type, stage Stage) *Session {
	return &Session{
		UserId:     userId,
		ActionType: action,
		Stage:      stage,
		Data:       make(map[SessionKey]interface{}),
	}
}

func (s *Session) AddData(key SessionKey, value interface{}) {
	s.Data[key] = value
}

func (s *Session) GetStringData(key SessionKey) string {
	return s.Data[key].(string)
}

func NewInMemorySessionModel() *SessionModel {
	return &SessionModel{session: make(map[int]*Session)}
}

func (m *SessionModel) FindSession(userId int) *Session {
	m.mux.Lock()
	defer m.mux.Unlock()

	ses, ok := m.session[userId]

	if !ok {
		return nil
	}

	return ses
}

func (m *SessionModel) ClearSession(userId int) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.session, userId)
}

func (m *SessionModel) UpdateSession(session *Session) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.session[session.UserId] = session
}
