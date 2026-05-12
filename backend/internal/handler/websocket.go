package handler

import (
	"log"
	"net/http"
	"sync"

	"cctv-backend/internal/stream"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

type WSHandler struct {
	manager *stream.Manager
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewWSHandler(manager *stream.Manager) *WSHandler {
	h := &WSHandler{
		manager: manager,
		clients: make(map[*websocket.Conn]bool),
	}
	go h.broadcastEvents()
	return h
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// Keep connection alive until client disconnects
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			break
		}
	}
}

func (h *WSHandler) broadcastEvents() {
	for event := range h.manager.Events() {
		h.mu.Lock()
		for client := range h.clients {
			if err := client.WriteJSON(event); err != nil {
				log.Printf("[ws] write error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}
