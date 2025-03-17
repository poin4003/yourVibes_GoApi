package socket_hub

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
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

	if oldConn, ok := hub.connections[userId]; ok {
		oldConn.Close()
	}

	hub.connections[userId] = conn
	log.Printf("WebSocket connection added for user_id: %s", userId)

	go hub.monitorConnection(userId, conn)
}

// remove connection to hub
func (hub *WebSocketHub) RemoveConnection(userId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if conn, ok := hub.connections[userId]; ok {
		conn.Close()
		delete(hub.connections, userId)
		log.Printf("WebSocket connection removed for user_id: %s", userId)
	}
}

// monitor connection
func (hub *WebSocketHub) monitorConnection(userId string, conn *websocket.Conn) {
	defer hub.RemoveConnection(userId)

	// Read message to keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket disconnected closed unexpected for user_id: %s, error: %v", userId, err)
			}
			return
		}
	}
}

// Send notification to User
func (hub *WebSocketHub) SendNotification(userId string, notification interface{}) error {
	hub.mu.RLock()
	conn, ok := hub.connections[userId]
	hub.mu.RUnlock()

	if !ok {
		log.Printf("No websocket connection found for user_id: %s", userId)
		return nil
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		hub.RemoveConnection(userId)
		return err
	}

	err := conn.WriteJSON(notification)
	if err != nil {
		log.Printf("Failed to send WebSocket notification to userId: %s, error: %v", userId, err)
		hub.RemoveConnection(userId)
		return err
	}

	log.Printf("WebSocket notification sent to userId: %s", userId)
	return nil
}
