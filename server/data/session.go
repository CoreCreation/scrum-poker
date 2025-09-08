package data

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Name string    `json:"name"`
	UUID uuid.UUID `json:"uuid"`
	Vote int64     `json:"vote"`
}

type Session struct {
	votesVisible bool
	voteOptions  string
	connections  map[*websocket.Conn]*Connection
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
	case "HideVotes":
		fmt.Println("Votes Hidden")
		s.votesVisible = false
	case "ShowVotes":
		fmt.Println("Votes Shown")
		s.votesVisible = true
	case "SetOptions":
		fmt.Println("Set Vote Options to:", msg.Body)
		s.voteOptions = msg.Body
	case "Init":
		// No-op just send state
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
	state := State{
		VotesVisible: s.votesVisible,
		VoteOptions:  s.voteOptions,
		UserData:     values,
	}

	for ws := range s.connections {
		go func(ws *websocket.Conn) {
			if err := ws.WriteJSON(&state); err != nil {
				fmt.Println("Write error:", err)
			}
		}(ws)
	}
}
