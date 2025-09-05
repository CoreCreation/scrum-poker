package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CoreCreation/scrum-poker/server/data"
	"github.com/google/uuid"
)

type SessionsHandler struct {
	data           *data.Data
	sessionHandler SessionHandler
}

func NewSessionsHandler(data *data.Data) *SessionsHandler {
	return &SessionsHandler{
		data:           data,
		sessionHandler: *NewSessionHandler(data),
	}
}

func (s *SessionsHandler) GetRouter() *http.ServeMux {
	sessionsRouter := http.NewServeMux()

	// Mount own routes
	sessionsRouter.HandleFunc("POST /sessions/create", s.CreateSession)
	sessionsRouter.HandleFunc("GET /sessions/{uuid}", s.GetSession)

	sessionRouter := s.sessionHandler.GetRouter()
	sessionRouter.Handle("/sessions/", http.StripPrefix("/sessions", sessionRouter))
	return sessionsRouter
}

type CreateSessionResponse struct {
	UUID string `json:"uuid"`
}

func (s *SessionsHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	uuid := s.data.Sessions.CreateSession()
	body := CreateSessionResponse{
		UUID: uuid.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		fmt.Println("Unable to encode json for uuid", err)
	}
	fmt.Println("Created Session UUID:", uuid)
}

func (s *SessionsHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		fmt.Println("Unable to parse session UUID:", uuidString, err)
		http.Error(w, "Bad UUID", http.StatusBadRequest)
		return
	}
	_, err = s.data.Sessions.GetSession(uuid)
	if err != nil {
		fmt.Println("Status check of Session:", uuidString)
		w.WriteHeader(http.StatusOK)
	}
}
