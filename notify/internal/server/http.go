package server

import (
	"context"
	"log"
	"net/http"

	"github.com/XRS0/ToTalkB/notify/internal/websocket"
)

type HTTPServer struct {
	wsManager *websocket.Manager
	server    *http.Server
	jwtKey    []byte
}

func NewHTTPServer(wsManager *websocket.Manager, jwtKey []byte) *HTTPServer {
	return &HTTPServer{
		wsManager: wsManager,
		jwtKey:    jwtKey,
	}
}

func (s *HTTPServer) Start(addr string) error {
	// Создаем WebSocket обработчик
	wsHandler := websocket.NewHandler(s.wsManager, s.jwtKey)

	// Регистрируем маршруты
	http.Handle("/ws", wsHandler)

	// Создаем HTTP сервер
	s.server = &http.Server{
		Addr: addr,
	}

	// Запускаем сервер
	log.Printf("Starting HTTP server on %s", addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
