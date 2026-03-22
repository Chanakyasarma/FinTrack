package handler

import (
	"log"
	"net/http"

	"fintrack/internal/websocket"
	"fintrack/pkg/middleware"
	"fintrack/pkg/respond"

	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := websocket.NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
