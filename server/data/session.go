package data

import (
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

type Client struct {
	Username    string `json:"username"`
	Vote        int64  `json:"vote"`
	Active      bool   `json:"active"`
	uuid        uuid.UUID
	connections map[*websocket.Conn]bool
	// cancel          context.CancelFunc
	mu              sync.Mutex
	mostRecentState atomic.Value
	havePing        atomic.Bool
}

func (c *Client) setData(cmd *ClientCommand) error {
	if cmd.Name != nil && *cmd.Name != c.Username {
		fmt.Println("Changing Username to", *cmd.Name)
		c.Username = *cmd.Name
	}
	if cmd.Vote != nil {
		newVote, err := cmd.Vote.Int64()
		if err != nil {
			fmt.Println("Unable to parse into int64", *cmd.Vote)
			return errors.New("Unable to parse number")
		}
		if newVote != c.Vote {
			fmt.Println("Vote Cast for", newVote)
			c.Vote = newVote
		}
	}
	if cmd.Active != nil && *cmd.Active != c.Active {
		fmt.Println("User changing active:", *cmd.Active)
		if *cmd.Active == false {
			c.Vote = -1
		}
		c.Active = *cmd.Active
	}
	return nil
}

type Session struct {
	parent              *Sessions
	votesVisible        bool
	voteOptions         string
	clients             map[uuid.UUID]*Client
	idleTimer           *time.Timer
	uuid                uuid.UUID
	toggleCooldown      atomic.Bool
	toggleCooldownTimer *time.Timer
}

func NewSession(parent *Sessions, sid uuid.UUID) *Session {
	session := &Session{
		parent:       parent,
		uuid:         sid,
		votesVisible: false,
		voteOptions:  "1, 2, 3, 5, 8, 13, 20",
		clients:      make(map[uuid.UUID]*Client),
	}

	session.idleTimer = time.AfterFunc(timeLimit, func() {
		if !session.NeedTimer() {
			return
		}

		session.parent.RemoveSession(session.uuid)
	})

	return session
}

func (s *Session) NeedTimer() bool {
	needTimer := true
	for _, client := range s.clients {
		needTimer = needTimer && len(client.connections) == 0
	}
	return needTimer
}

func (s *Session) handleLifeTimer() {
	if s.NeedTimer() {
		fmt.Println("All clients disconnected, starting session timer")
		s.idleTimer.Reset(timeLimit)
	} else {
		fmt.Println("Some clients connected, stopping session timer")
		s.idleTimer.Stop()
	}
}

func (s *Session) HandleConnection(clientId uuid.UUID, connection *websocket.Conn) {

	client, ok := s.clients[clientId]
	// Client already exists, add its connection
	if ok {

		// Cancel the current connection
		// client.cancel()

		// and then create new context
		// ctx, cancel := context.WithCancel(context.Background())
		// go func() {
		// 	<-ctx.Done()
		// 	connection.SetReadDeadline(time.Now())
		// }()
		// client.cancel = cancel

		// and then set new connection
		client.connections[connection] = true

		// trigger fresh data to be sent
		if len(client.connections) > 1 {
			s.sendState(client)
		} else {
			s.broadcastState()
		}

	} else {
		// ctx, cancel := context.WithCancel(context.Background())
		client = &Client{
			uuid:     uuid.New(),
			Username: "",
			Vote:     -1,
			Active:   true,
			// cancel:      cancel,
			connections: map[*websocket.Conn]bool{connection: true},
		}

		s.clients[clientId] = client

		// If context is canceled, kill the connection right away
		// go func() {
		// 	<-ctx.Done()
		// 	connection.SetReadDeadline(time.Now())
		// }()

		s.sendInit(client)
	}

	s.readLoop(connection, s.clients[clientId])
}

func (s *Session) RemoveConnection(cid uuid.UUID, connection *websocket.Conn) {
	client, ok := s.clients[cid]
	if !ok {
		fmt.Println("Unable to get client", cid)
	}
	delete(client.connections, connection)
	s.broadcastState()
	s.handleLifeTimer()
}

type ClientCommand struct {
	Type   *string      `json:"type"`
	Body   *string      `json:"body"`
	Name   *string      `json:"username"`
	Vote   *json.Number `json:"vote"`
	Active *bool        `json:"active"`
}

type ServerCommand struct {
	Type string `json:"type"`
}

func (s *Session) readLoop(connection *websocket.Conn, client *Client) {
	for {
		var command ClientCommand
		if err := connection.ReadJSON(&command); err != nil {
			fmt.Println("Unable to read Client Command:", err)
			break
		}
		s.handleMessage(&command, client)
	}
}

func (s *Session) handleMessage(cmd *ClientCommand, client *Client) {
	switch *cmd.Type {
	case "UpdateData":
		err := client.setData(cmd)
		if err != nil {
			break
		}
	case "ClearVotes":
		if s.toggleCooldown.Load() {
			fmt.Println("Toggle is on cooldown, skipping")
			break
		}
		fmt.Println("All Votes Cleared")
		for _, client := range s.clients {
			client.Vote = -1
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
		fmt.Println("Set Vote Options to:", cmd.Body)
		s.voteOptions = *cmd.Body
	}

	s.broadcastState()
}

func (s *Session) setTimer() {
	s.toggleCooldown.Store(true)
	s.toggleCooldownTimer = time.AfterFunc(1*time.Second, func() {
		s.toggleCooldown.Swap(false)
	})
}

type State struct {
	VotesVisible bool      `json:"votesVisible"`
	VoteOptions  string    `json:"voteOptions"`
	ClientData   []*Client `json:"clientData"`
	Username     string    `json:"username"`
	Active       bool      `json:"active"`
}

func (s *Session) collectionClientData() []*Client {
	values := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		if len(client.connections) > 0 {
			values = append(values, client)
		}
	}
	return values
}

func (s *Session) createState(client *Client, clientData []*Client) ([]byte, error) {
	jd, err := json.Marshal(State{
		Username:     client.Username,
		VotesVisible: s.votesVisible,
		VoteOptions:  s.voteOptions,
		Active:       client.Active,
		ClientData:   clientData,
	})
	if err != nil {
		fmt.Println("Unable to marshal JSON")
		return nil, err
	}
	return jd, nil
}

func (s *Session) broadcastState() {
	clientData := s.collectionClientData()

	for _, client := range s.clients {
		jd, err := s.createState(client, clientData)
		if err != nil {
			return
		}
		go s.sendQueuedData(client, jd, false)
	}
}

func (s *Session) sendState(client *Client) {
	clientData := s.collectionClientData()

	jd, err := s.createState(client, clientData)
	if err != nil {
		return
	}
	go s.sendQueuedData(client, jd, false)
}

func (s *Session) sendInit(client *Client) {

	fmt.Println("Sending Init for a new client")

	jd, err := s.createServerCommand(&ServerCommand{
		Type: "Init",
	})
	if err != nil {
		return
	}
	go s.sendQueuedData(client, jd, false)
}

func (s *Session) createServerCommand(cmd *ServerCommand) ([]byte, error) {
	jd, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println("Unable to marshal JSON")
		return nil, err
	}
	return jd, nil
}

func (s *Session) sendPings() {
	for _, client := range s.clients {
		go s.sendQueuedData(client, nil, true)
	}
}

func (s *Session) sendQueuedData(client *Client, json []byte, ping bool) {
	fmt.Println("- Going to broadcast")
	if len(client.connections) == 0 {
		fmt.Println("Client connections == 0, skipping message")
		return
	}
	if json != nil {
		client.mostRecentState.Store(json)
	}
	if ping {
		client.havePing.Store(true)
	}
	if client.mu.TryLock() {
		fmt.Println("-- Got the lock")
		defer client.mu.Unlock()
		jd := client.mostRecentState.Swap([]byte(""))
		if jd != nil && len(jd.([]byte)) != 0 {
			fmt.Println("-- JSON Message Found")
			for connection := range client.connections {
				if err := connection.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
					fmt.Println("--- Write error while sending JSON:", err)
					return
				}
			}
			fmt.Println("-- Sent JSON")
		} else if client.havePing.Swap(false) {
			fmt.Println("-- New Ping, sending Ping")
			for connection := range client.connections {
				if err := connection.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Println("--- Write error while sending Ping:", err)
					return
				}
			}
		}
		for jd = client.mostRecentState.Swap([]byte("")); jd != nil && len(jd.([]byte)) != 0; jd = client.mostRecentState.Swap([]byte("")) {
			fmt.Println("-- New Data, sending JSON")
			for connection := range client.connections {
				if err := connection.WriteMessage(websocket.TextMessage, jd.([]byte)); err != nil {
					fmt.Println("--- Write error while sending JSON:", err)
					return
				}
			}
		}
	} else {
		fmt.Println("-- Unable to get the lock")
		return
	}
}
