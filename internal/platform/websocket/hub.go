package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single websocket client.
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID uint
	Roles  []string
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("🔌 WS Client registered: user_id=%d", client.UserID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("🔌 WS Client unregistered: user_id=%d", client.UserID)
		}
	}
}

// WSMessage matches the format expected by the frontend's useWebSocket hook.
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(msgType string, payload interface{}) {
	data, err := json.Marshal(WSMessage{Type: msgType, Payload: payload})
	if err != nil {
		log.Printf("⚠️ WS Hub: failed to marshal message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
			log.Printf("⚠️ WS Hub: client send channel full, closing connection")
			go h.Unregister(client)
		}
	}
}

// Unregister safely unregisters a client from the hub.
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
	c.Conn.Close()
}

func (c *Client) WritePump() {
	defer func() {
		c.Hub.Unregister(c)
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		n := len(c.Send)
		for i := 0; i < n; i++ {
			w.Write([]byte{'\n'})
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
	}()
	c.Conn.SetReadLimit(512)
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
