package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Sessions struct {
	Sessions map[uuid.UUID]*Session
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSessions() *Sessions {
	return &Sessions{
		Sessions: make(map[uuid.UUID]*Session),
	}
}

func (s *Sessions) CreateSession() uuid.UUID {
	uuid := uuid.New()
	s.Sessions[uuid] = NewSession(s, uuid)
	if len(s.Sessions) > 0 && s.ctx == nil {
		fmt.Println("Sessions added, creating heartbeat", uuid)
		s.ctx, s.cancel = context.WithCancel(context.Background())
		go s.heartbeat(s.ctx)
	}
	return uuid
}

func (s *Sessions) RemoveSession(uuid uuid.UUID) {
	fmt.Println("Timer is removing session:", uuid)
	delete(s.Sessions, uuid)
	if len(s.Sessions) == 0 {
		fmt.Println("All sessions removed, canceling heartbeat", uuid)
		s.cancel()
		s.ctx = nil
		s.cancel = nil
	}
}

func (s *Sessions) GetSession(uuid uuid.UUID) (*Session, error) {
	session, ok := s.Sessions[uuid]
	if !ok {
		return nil, errors.New("Session Not Found")
	}
	return session, nil
}

func (s *Sessions) heartbeat(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			for _, session := range s.Sessions {
				session.sendPings()
			}
		}
	}
}
