package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/XRS0/ToTalkB/auth/pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене здесь должна быть проверка origin
	},
}

// Handler handles WebSocket connections
type Handler struct {
	manager *Manager
	jwtKey  []byte
}

func NewHandler(manager *Manager, jwtKey []byte) *Handler {
	return &Handler{
		manager: manager,
		jwtKey:  jwtKey,
	}
}

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Получаем JWT токен из query параметров
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	// Парсим и проверяем токен
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Создаем объект пользователя из данных токена
	user := &pkg.User{
		Id:    claims.UserID,
		Login: claims.Login,
		Name:  claims.Name,
		Role:  claims.Role,
	}

	// Обновляем HTTP соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	wsUser := &WebSocketUser{
		User:    user,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Manager: h.manager,
	}

	// Регистрируем пользователя
	h.manager.register <- wsUser

	// Запускаем горутины для чтения и записи
	go wsUser.writePump()
	go wsUser.readPump()
}

func (u *WebSocketUser) readPump() {
	defer func() {
		u.Manager.unregister <- u
		u.Conn.Close()
	}()

	u.Conn.SetReadLimit(512)
	u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := u.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
	}
}

func (u *WebSocketUser) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-u.Send:
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
