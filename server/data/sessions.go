package data

import (
	"errors"

	"github.com/google/uuid"
)

type Sessions struct {
	Sessions map[uuid.UUID]*Session
}

func NewSessions() *Sessions {
	return &Sessions{
		Sessions: make(map[uuid.UUID]*Session),
	}
}

func (s *Sessions) CreateSession() uuid.UUID {
	uuid := uuid.New()
	s.Sessions[uuid] = NewSession()
	return uuid
}

func (s *Sessions) GetSession(uuid uuid.UUID) (*Session, error) {
	session, ok := s.Sessions[uuid]
	if !ok {
		return nil, errors.New("Session Not Found")
	}
	return session, nil
}
