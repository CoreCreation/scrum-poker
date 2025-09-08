package data

import (
	"errors"
	"fmt"

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
	s.Sessions[uuid] = NewSession(s, uuid)
	return uuid
}

func (s *Sessions) RemoveSession(uuid uuid.UUID) {
	fmt.Println("Timer is removing session:", uuid)
	delete(s.Sessions, uuid)
}

func (s *Sessions) GetSession(uuid uuid.UUID) (*Session, error) {
	session, ok := s.Sessions[uuid]
	if !ok {
		return nil, errors.New("Session Not Found")
	}
	return session, nil
}
