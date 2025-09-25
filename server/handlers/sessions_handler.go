package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CoreCreation/scrum-poker/server/data"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type SessionsHandler struct {
	data *data.Sessions
}

func NewSessionsHandler(data *data.Sessions) *SessionsHandler {
	return &SessionsHandler{
		data: data,
	}
}

func (s *SessionsHandler) GetRouter() *http.ServeMux {
	sessionsRouter := http.NewServeMux()

	// Mount own routes
	sessionsRouter.HandleFunc("POST /sessions/create", s.CreateSession)
	sessionsRouter.HandleFunc("GET /sessions/{uuid}", s.GetSession)
	sessionsRouter.HandleFunc("/sessions/{sid}/join/{cid}", s.JoinSession)

	return sessionsRouter
}

type CreateSessionResponse struct {
	UUID string `json:"uuid"`
}

func (s *SessionsHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	uuid := s.data.CreateSession()
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
	_, err = s.data.GetSession(uuid)
	if err != nil {
		fmt.Println("Error getting session", err)
		http.Error(w, "Error Getting Session", http.StatusInternalServerError)
		return
	}
	fmt.Println("Status check of Session:", uuidString)
	w.WriteHeader(http.StatusOK)
}

func (s *SessionsHandler) JoinSession(w http.ResponseWriter, r *http.Request) {

	// Handle Session
	sidString := r.PathValue("sid")
	sid, err := uuid.Parse(sidString)
	if err != nil {
		fmt.Println("Unable to parse session UUID:", sidString, err)
		http.Error(w, "Bad Session ID", http.StatusBadRequest)
		return
	}
	fmt.Println("Client trying to open WebSocket Connection to Session:", sidString)
	session, err := s.data.GetSession(sid)
	if err != nil {
		fmt.Println("Error getting session", err)
		http.Error(w, "Error Getting Session", http.StatusInternalServerError)
		return
	}
	fmt.Println("Request to join Session:", sidString)

	// Get Client ID
	cidString := r.PathValue("cid")
	cid, err := uuid.Parse(cidString)
	if err != nil {
		fmt.Println("Unable to parse session UUID:", cidString, err)
		http.Error(w, "Bad Client ID", http.StatusBadRequest)
		return
	}

	// Upgrade to get Websocket connection
	fmt.Println("Upgrading to WebSocket for Client:", cidString)
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error occured while trying to upgrade connection to WebSocket:", err)
		return
	}
	defer func() {
		session.RemoveConnection(cid, connection)
		connection.Close()
	}()
	connection.SetReadLimit(1 << 20)
	connection.SetReadDeadline(time.Now().Add(20 * time.Second))
	connection.SetPongHandler(func(string) error {
		connection.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})

	// Handle client (new/old)
	fmt.Println("Client Connecting:", cidString)
	session.HandleConnection(cid, connection)
}
