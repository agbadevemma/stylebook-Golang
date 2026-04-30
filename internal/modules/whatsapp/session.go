package whatsapp

import "github.com/stylebookgo/agbadevemma/internal/utils"

type Session struct {
	Phone string
	State string
	Data  map[string]interface{}
}

type SessionService struct {
	store map[string]*Session
}

func NewSessionService() *SessionService {
	return &SessionService{
		store: make(map[string]*Session),
	}
}

func (s *SessionService) GetOrCreate(phone string) *Session {
	if session, ok := s.store[phone]; ok {
		return session
	}

	session := &Session{
		Phone: phone,
		State: string(utils.StateStart),
		Data:  make(map[string]any),
	}

	s.store[phone] = session
	return session
}

func (s *SessionService) Save(session *Session) {
	s.store[session.Phone] = session
}
