package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var timeLimit time.Duration = 10 * time.Minute

type Connection struct {
	Name            string    `json:"name"`
	UUID            uuid.UUID `json:"uuid"`
	Vote            int64     `json:"vote"`
	Active          bool      `json:"active"`
	mu              sync.Mutex
	mostRecentState atomic.Value
	havePing        atomic.Bool
	cancel          context.CancelFunc
}

type Session struct {
	parent              *Sessions
	votesVisible        bool
	voteOptions         string
	connections         map[*websocket.Conn]*Connection
	idleTimer           *time.Timer
	uuid                uuid.UUID
	toggleCooldown      atomic.Bool
	toggleCooldownTimer *time.Timer
}

func NewSession(parent *Sessions, uuid uuid.UUID) *Session {
	session := &Session{
		parent:       parent,
		uuid:         uuid,
		votesVisible: false,
		voteOptions:  "1, 2, 3, 5, 8, 12",
		connections:  make(map[*websocket.Conn]*Connection),
	}

	session.idleTimer = time.AfterFunc(timeLimit, func() {
		if len(session.connections) > 0 {
			return
		}

		session.parent.RemoveSession(session.uuid)
	})

	return session
}

func (s *Session) AddConnection(connection *websocket.Conn) {
	fmt.Println("Adding a connection and stopping session timer")
	s.idleTimer.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	s.connections[connection] = &Connection{
		UUID:   uuid.New(),
		Name:   "",
		Vote:   -1,
		Active: true,
		cancel: cancel,
	}

	// If context is canceled, kill the connection right away
	go func() {
		<-ctx.Done()
		connection.SetReadDeadline(time.Now())
	}()

	s.readLoop(connection, s.connections[connection])
}

func (s *Session) RemoveConnection(connection *websocket.Conn) {
	delete(s.connections, connection)
	s.sendState()
	if len(s.connections) == 0 {
		fmt.Println("All connections removed, starting session timer")
		s.idleTimer.Reset(timeLimit)
	}
}

type ClientCommand struct {
	Type   *string      `json:"type"`
	Body   *string      `json:"body"`
	Name   *string      `json:"username"`
	Vote   *json.Number `json:"vote"`
	Active *bool        `json:"active"`
}

func (s *Session) readLoop(connection *websocket.Conn, data *Connection) {
	for {
		var command ClientCommand
		if err := connection.ReadJSON(&command); err != nil {
			fmt.Println("Unable to read Client Command:", err)
			break
		}
		s.handleMessage(command, data)
	}
}

func setData(msg ClientCommand, data *Connection) error {
	if msg.Name != nil && *msg.Name != data.Name {
		fmt.Println("Changing Username to", msg.Body)
		data.Name = *msg.Name
	}
	if msg.Vote != nil {
		newVote, err := msg.Vote.Int64()
		if err != nil {
			fmt.Println("Unable to parse into int64", msg.Body)
			return errors.New("Unable to parse number")
		}
		if newVote != data.Vote {
			fmt.Println("Vote Cast for", newVote)
			data.Vote = newVote
		}
	}
	if msg.Active != nil && *msg.Active != data.Active {
		fmt.Println("User changing active:", msg.Active)
		if *msg.Active == false {
			data.Vote = -1
		}
		data.Active = *msg.Active
	}
	return nil
}

func (s *Session) handleMessage(msg ClientCommand, data *Connection) {
	switch *msg.Type {
	case "UpdateData":
		err := setData(msg, data)
		if err != nil {
			break
		}
	case "ClearVotes":
		if s.toggleCooldown.Load() {
			fmt.Println("Toggle is on cooldown, skipping")
			break
		}
		fmt.Println("All Votes Cleared")
		for _, data := range s.connections {
			data.Vote = -1
		}
		s.votesVisible = false
		s.setTimer()
	case "ShowVotes":
		if s.toggleCooldown.Load() {
			fmt.Println("Toggle is on cooldown, skipping")
			break
		}
		fmt.Println("Votes Shown")
		s.votesVisible = true
		s.setTimer()
	case "SetOptions":
		fmt.Println("Set Vote Options to:", msg.Body)
		s.voteOptions = *msg.Body
	case "Init":
		fmt.Println("Init")
		if msg.Body != nil && len(*msg.Body) > 0 {
			fmt.Println("Connection made with already existent UUID, evicting old connection")
			parsed, err := uuid.Parse(*msg.Body)
			if err != nil {
				fmt.Println("UUID passed for connection can not be parsed, ignoring", err)
				break
			}
			for _, connection := range s.connections {
				if connection.UUID == parsed {
					fmt.Println("Connection found, canceling")
					connection.cancel()
				}
			}
		}
		err := setData(msg, data)
		if err != nil {
			fmt.Println("Unable to set data during init", err)
			break
		}
	}

	s.sendState()
}

func (s *Session) setTimer() {
	s.toggleCooldown.Store(true)
	s.toggleCooldownTimer = time.AfterFunc(1*time.Second, func() {
		s.toggleCooldown.Swap(false)
	})
}

type State struct {
	UserID       uuid.UUID     `json:"userId"`
	VotesVisible bool          `json:"votesVisible"`
	VoteOptions  string        `json:"voteOptions"`
	UserData     []*Connection `json:"userData"`
	Username     string        `json:"username"`
}

func (s *Session) sendState() {
	values := make([]*Connection, 0, len(s.connections))
	for _, v := range s.connections {
		values = append(values, v)
	}

	for ws, data := range s.connections {
		jd, err := json.Marshal(State{
			UserID:       data.UUID,
			Username:     data.Name,
			VotesVisible: s.votesVisible,
			VoteOptions:  s.voteOptions,
			UserData:     values,
		})
		if err != nil {
			fmt.Println("Unable to marshal JSON")
			return
		}
		go s.sendQueuedData(ws, data, jd, false)
	}
}

func (s *Session) sendPings() {
	for ws, data := range s.connections {
		go s.sendQueuedData(ws, data, nil, true)
	}
}

func (s *Session) sendQueuedData(ws *websocket.Conn, data *Connection, json []byte, ping bool) {
	fmt.Println("- Going to broadcast")
	if json != nil {
		data.mostRecentState.Store(json)
	}
	if ping {
		data.havePing.Store(true)
	}
	if data.mu.TryLock() {
		fmt.Println("-- Got the lock")
		defer data.mu.Unlock()
		jd := data.mostRecentState.Swap([]byte(""))
		if jd != nil && len(jd.([]byte)) != 0 {
			fmt.Println("-- JSON Message Found")
			if err := ws.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
				fmt.Println("--- Write error while sending JSON:", err)
				return
			}
			fmt.Println("-- Sent JSON")
		} else if data.havePing.Swap(false) {
			fmt.Println("-- New Ping, sending Ping")
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("--- Write error while sending Ping:", err)
				return
			}
		}
		for jd = data.mostRecentState.Swap([]byte("")); jd != nil && len(jd.([]byte)) != 0; jd = data.mostRecentState.Swap([]byte("")) {
			fmt.Println("-- New Data, sending JSON")
			if err := ws.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
				fmt.Println("--- Write error while sending JSON:", err)
				return
			}
		}
	} else {
		fmt.Println("-- Unable to get the lock")
		return
	}
}
