package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type EventType string

const (
	EventTransactionCreated EventType = "transaction.created"
	EventBalanceUpdated     EventType = "balance.updated"
)

type Event struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	mu     sync.Mutex
}

type Hub struct {
	clients    map[string]map[*Client]struct{} // userID → clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan userMessage
	mu         sync.RWMutex
}

type userMessage struct {
	userID  string
	payload []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan userMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]struct{})
			}
			h.clients[client.userID][client] = struct{}{}
			h.mu.Unlock()
			log.Printf("ws: client registered for user %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.userID)
				}
			}
			h.mu.Unlock()
			close(client.send)
			log.Printf("ws: client unregistered for user %s", client.userID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.userID]
			h.mu.RUnlock()
			for c := range clients {
				select {
				case c.send <- msg.payload:
				default:
					close(c.send)
					h.mu.Lock()
					delete(h.clients[msg.userID], c)
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) SendToUser(userID string, event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("ws: failed to marshal event: %v", err)
		return
	}
	h.broadcast <- userMessage{userID: userID, payload: payload}
}

func (c *Client) WritePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for msg := range c.send {
		c.mu.Lock()
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		c.mu.Unlock()
		if err != nil {
			break
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
