package socket_hub

import (
	"github.com/gorilla/websocket"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"sync"
)

type WebSocketHub struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		connections: make(map[string]*websocket.Conn),
	}
}

// Add connection to hub
func (hub *WebSocketHub) AddConnection(userId string, conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.connections[userId] = conn
}

// Remove connection to hub
func (hub *WebSocketHub) RemoveConnection(userId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if conn, ok := hub.connections[userId]; ok {
		conn.Close()
		delete(hub.connections, userId)
	}
}

// Send notification to User
func (hub *WebSocketHub) SendNotification(userId string, notification *consts.NotificationSocketResponse) error {
	hub.mu.RLock()
	conn, ok := hub.connections[userId]
	hub.mu.RUnlock()

	if !ok {
		return nil
	}

	err := conn.WriteJSON(notification)
	if err != nil {
		hub.RemoveConnection(userId)
		return err
	}

	return nil
}
