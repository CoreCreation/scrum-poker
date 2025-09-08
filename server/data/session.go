package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Name      string    `json:"name"`
	UUID      uuid.UUID `json:"uuid"`
	Vote      int64     `json:"vote"`
	Active    bool      `json:"active"`
	mu        sync.Mutex
	outOfDate atomic.Bool
}

type Session struct {
	votesVisible    bool
	voteOptions     string
	connections     map[*websocket.Conn]*Connection
	mostRecentState atomic.Value
}

func NewSession() *Session {
	return &Session{
		votesVisible: false,
		voteOptions:  "1, 2, 3, 5, 8",
		connections:  make(map[*websocket.Conn]*Connection),
	}
}

func (s *Session) AddConnection(connection *websocket.Conn) {
	s.connections[connection] = &Connection{
		UUID: uuid.New(),
		Name: "No Username",
		Vote: -1,
	}
	s.readLoop(connection, s.connections[connection])
}

type ClientCommand struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

func (s *Session) readLoop(connection *websocket.Conn, data *Connection) {
	for {
		var command ClientCommand

		if err := connection.ReadJSON(&command); err != nil {
			fmt.Println("Unable to read Client Command:", err)
			delete(s.connections, connection)
			s.sendState()
			break
		}
		s.handleMessage(command, data)
	}
}

func (s *Session) handleMessage(msg ClientCommand, data *Connection) {
	switch msg.Type {
	case "SetName":
		fmt.Println("Changing Username to", msg.Body)
		data.Name = msg.Body
	case "CastVote":
		vote, err := strconv.ParseInt(msg.Body, 10, 64)
		if err != nil {
			fmt.Println("Unable to parse into int64", msg.Body)
			break
		}
		fmt.Println("Vote Cast for", msg.Body)
		data.Vote = vote
	case "ClearVotes":
		fmt.Println("All Votes Cleared")
		for _, data := range s.connections {
			data.Vote = -1
		}
		s.votesVisible = false
	case "ShowVotes":
		fmt.Println("Votes Shown")
		s.votesVisible = true
	case "SetOptions":
		fmt.Println("Set Vote Options to:", msg.Body)
		s.voteOptions = msg.Body
	case "LeaveVote":
		fmt.Println("User leaving vote")
		data.Active = false
	case "JoinVote":
		fmt.Println("User joining vote")
		data.Active = true
	case "Init":
		fmt.Println("Init")
	}

	s.sendState()
}

type State struct {
	VotesVisible bool          `json:"votesVisible"`
	VoteOptions  string        `json:"voteOptions"`
	UserData     []*Connection `json:"userData"`
}

func (s *Session) sendState() {
	values := make([]*Connection, 0, len(s.connections))
	for _, v := range s.connections {
		values = append(values, v)
	}

	jd, err := json.Marshal(State{
		VotesVisible: s.votesVisible,
		VoteOptions:  s.voteOptions,
		UserData:     values,
	})
	if err != nil {
		fmt.Println("Unable to marshal JSON")
		return
	}
	s.mostRecentState.Store(jd)

	for ws, data := range s.connections {
		go s.sendData(ws, data)
	}
}

func (s *Session) sendData(ws *websocket.Conn, data *Connection) {
	fmt.Println("Going to broadcast data", data.outOfDate.Load())
	if data.mu.TryLock() {
		fmt.Println("Got the lock")
		defer data.mu.Unlock()
		jd := s.mostRecentState.Load()
		if jd == nil {
			fmt.Println("Cached JSON not found")
			return
		}
		if err := ws.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
			fmt.Println("Write error:", err)
			data.outOfDate.Store(false)
			return
		}
		fmt.Println("Sent JSON", data.outOfDate.Load())
		for data.outOfDate.Swap(false) {
			fmt.Println("After sending data, outOfDate was true, sending data again")
			jd := s.mostRecentState.Load()
			if jd == nil {
				fmt.Println("Cached JSON not found")
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
				fmt.Println("Write error:", err)
				data.outOfDate.Store(false)
				return
			}
		}
	} else {
		data.outOfDate.Store(true)
		fmt.Println("Unable to get lock, will send data again")
		return
	}
}
