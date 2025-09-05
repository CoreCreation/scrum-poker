package handlers

import (
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
)

type SessionHandler struct {
	data *data.Data
}

func NewSessionHandler(data *data.Data) *SessionHandler {
	return &SessionHandler{
		data: data,
	}
}

func (s *SessionHandler) GetRouter() *http.ServeMux {
	sessionRouter := http.NewServeMux()

	return sessionRouter
}
