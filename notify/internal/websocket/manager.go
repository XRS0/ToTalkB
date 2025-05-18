package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/XRS0/ToTalkB/auth/pkg"
	"github.com/gorilla/websocket"
)

// WebSocketUser represents a connected WebSocket user
type WebSocketUser struct {
	User    *pkg.User
	Conn    *websocket.Conn
	Send    chan []byte
	Manager *Manager
	mu      sync.Mutex
}

// Manager handles all WebSocket connections
type Manager struct {
	users      map[int]*WebSocketUser // map[userID]*WebSocketUser
	register   chan *WebSocketUser
	unregister chan *WebSocketUser
	mu         sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewManager() *Manager {
	return &Manager{
		users:      make(map[int]*WebSocketUser),
		register:   make(chan *WebSocketUser),
		unregister: make(chan *WebSocketUser),
	}
}

func (m *Manager) Start() {
	for {
		select {
		case user := <-m.register:
			m.mu.Lock()
			m.users[user.User.Id] = user
			m.mu.Unlock()
			log.Printf("User registered: %s", user.User.Login)

		case user := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.users[user.User.Id]; ok {
				delete(m.users, user.User.Id)
				close(user.Send)
			}
			m.mu.Unlock()
			log.Printf("User unregistered: %s", user.User.Login)
		}
	}
}

// SendToUser sends a message to a specific user
func (m *Manager) SendToUser(userID int, messageType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := Message{
		Type:    messageType,
		Payload: payloadBytes,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m.mu.RLock()
	user, exists := m.users[userID]
	m.mu.RUnlock()

	if !exists {
		return nil // User is not connected
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	return user.Conn.WriteMessage(websocket.TextMessage, messageBytes)
}
