package websocket

import (
	"log"
	"net/http"

	"university-erp-backend/internal/platform/auth"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for the API server
	},
}

// Handler handles websocket upgrade requests.
type Handler struct {
	hub    *Hub
	jwtMgr *auth.JWTManager
}

func NewHandler(hub *Hub, jwtMgr *auth.JWTManager) *Handler {
	return &Handler{
		hub:    hub,
		jwtMgr: jwtMgr,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Println("⚠️ WS Upgrade: missing token")
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}

	claims, err := h.jwtMgr.ValidateToken(token)
	if err != nil {
		log.Printf("⚠️ WS Upgrade: invalid token: %v", err)
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("⚠️ WS Upgrade: failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		Hub:    h.hub,
		Conn:   conn,
		UserID: claims.UserID,
		Roles:  claims.Roles,
		Send:   make(chan []byte, 256),
	}
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
