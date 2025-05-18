package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/XRS0/ToTalkB/chat/pkg"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	db     *sqlx.DB
	userId string
	chatId string
}

type IncomingMessage struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

type OutgoingMessage struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
	Time    string `json:"time"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var incomingMsg IncomingMessage
		if err := json.Unmarshal(msgBytes, &incomingMsg); err != nil {
			log.Printf("invalid json: %v", err)
			continue
		}
		log.Println(incomingMsg)

		content := strings.TrimSpace(incomingMsg.Message)
		now := time.Now()

		_, err = c.db.Exec(
			"INSERT INTO messages (chat_id, sender_id, created_at, content) VALUES ($1, $2, $3, $4)",
			c.chatId,
			c.userId,
			now,
			content,
		)
		if err != nil {
			log.Printf("Failed to save message to DB: %v", err)
			continue
		}

		outMsg := OutgoingMessage{
			Content: content,
			Sender:  incomingMsg.Sender,
			Time:    now.Format("15:04"),
		}

		jsonMsg, err := json.Marshal(outMsg)
		if err != nil {
			log.Printf("Failed to marshal json: %v", err)
			continue
		}

		c.hub.broadcast <- jsonMsg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for range n {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, db *sqlx.DB, userId, chatId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), db: db, userId: userId, chatId: chatId}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()

	messages, err := loadChatHistory(db, chatId)
	if err != nil {
		log.Printf("Failed to load chat history: %v", err)
	} else {
		for _, msg := range messages {
			outMsg := OutgoingMessage{
				Content: msg.Content,
				Sender:  msg.Sender,
				Time:    msg.CreatedAt.Format("15:04"),
			}
			log.Println(outMsg)
			jsonMsg, err := json.Marshal(outMsg)
			if err != nil {
				log.Printf("Failed to marshal history message: %v", err)
				continue
			}
			client.send <- jsonMsg
		}
	}
}

func loadChatHistory(db *sqlx.DB, chatId string) ([]pkg.Message, error) {
	var messages []pkg.Message
	query := `
        SELECT m.content, u.name AS sender, m.created_at
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.chat_id = $1
        ORDER BY m.created_at ASC
    `
	err := db.Select(&messages, query, chatId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
