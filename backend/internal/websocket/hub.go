package websocket

import (
"ai-desktop-assistant/backend/internal/events"
"ai-desktop-assistant/backend/internal/logging"
"sync"
)

// Client represents a WebSocket client
type Client struct {
ID     string
Send   chan []byte
Hub    *Hub
}

// Hub maintains active WebSocket clients
type Hub struct {
clients    map[*Client]bool
broadcast  chan []byte
register   chan *Client
unregister chan *Client
eventBus   *events.Bus
logger     *logging.Logger
mu         sync.RWMutex
running    bool
}

// NewHub creates a new WebSocket hub
func NewHub(eventBus *events.Bus, logger *logging.Logger) *Hub {
return &Hub{
clients:    make(map[*Client]bool),
broadcast:  make(chan []byte, 256),
register:   make(chan *Client),
unregister: make(chan *Client),
eventBus:   eventBus,
logger:     logger,
running:    false,
}
}

// Start begins the hub's event loop
func (h *Hub) Start() {
h.mu.Lock()
h.running = true
h.mu.Unlock()

go h.run()
h.logger.Info("WebSocket hub started")
}

// Stop halts the hub
func (h *Hub) Stop() {
h.mu.Lock()
h.running = false
h.mu.Unlock()

close(h.broadcast)
h.logger.Info("WebSocket hub stopped")
}

func (h *Hub) run() {
for h.running {
select {
case client := <-h.register:
h.mu.Lock()
h.clients[client] = true
h.mu.Unlock()
h.logger.Info("Client connected", "id", client.ID)

case client := <-h.unregister:
h.mu.Lock()
if _, ok := h.clients[client]; ok {
delete(h.clients, client)
close(client.Send)
}
h.mu.Unlock()
h.logger.Info("Client disconnected", "id", client.ID)

case message := <-h.broadcast:
h.mu.RLock()
for client := range h.clients {
select {
case client.Send <- message:
default:
close(client.Send)
delete(h.clients, client)
}
}
h.mu.RUnlock()
}
}
}

// Register adds a new client
func (h *Hub) Register(client *Client) {
h.register <- client
}

// Unregister removes a client
func (h *Hub) Unregister(client *Client) {
h.unregister <- client
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(message []byte) {
h.broadcast <- message
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
h.mu.RLock()
defer h.mu.RUnlock()
return len(h.clients)
}
