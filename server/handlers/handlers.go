package handlers

import (
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
)

type Handlers struct {
	data *data.Data

	sessionsHandler *SessionsHandler
}

func NewHandlers(data *data.Data) *Handlers {
	return &Handlers{
		data:            data,
		sessionsHandler: NewSessionsHandler(data),
	}
}

func (h *Handlers) GetRouter() *http.ServeMux {
	apiRouter := http.NewServeMux()

	sessionsRouter := h.sessionsHandler.GetRouter()
	apiRouter.Handle("/api/", http.StripPrefix("/api", sessionsRouter))
	return apiRouter
}
