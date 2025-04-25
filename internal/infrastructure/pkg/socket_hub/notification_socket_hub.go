package socket_hub

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type NotificationSocketHub struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

func NewNotificationSocketHub() *NotificationSocketHub {
	return &NotificationSocketHub{
		connections: make(map[string]*websocket.Conn),
	}
}

// Add connection to hub
func (hub *NotificationSocketHub) AddConnection(userId string, conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if oldConn, ok := hub.connections[userId]; ok {
		oldConn.Close()
	}

	hub.connections[userId] = conn
	log.Printf("WebSocket connection added for user_id: %s", userId)
}

// remove connection to hub
func (hub *NotificationSocketHub) RemoveConnection(userId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	conn, ok := hub.connections[userId]
	if !ok {
		return
	}

	delete(hub.connections, userId)
	if conn != nil {
		conn.Close()
	}
	log.Printf("WebSocket connection removed for user_id: %s", userId)
}

// Send notification to User
func (hub *NotificationSocketHub) SendNotification(userId string, notification interface{}) error {
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
