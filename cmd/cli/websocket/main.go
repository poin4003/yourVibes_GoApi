package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Notification struct {
	From    User   `json:"from"`
	To      User   `json:"to"`
	Message string `json:"message"`
}

const (
	ConsumerGroup      = "notification-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

var ErrNoMessageFound = errors.New("message not found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessageFound
	}

	return userID, nil
}

type UserNotification map[string][]Notification

type NotificationStore struct {
	data UserNotification
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string, notification Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.data[userID]
}

func readMessages(ctx context.Context, store *NotificationStore, ws *websocket.Conn, done chan struct{}) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{KafkaServerAddress},
		Topic:   ConsumerTopic,
		GroupID: ConsumerGroup,
	})

	for {
		select {
		case <-done:
			log.Println("Stopping message reading...")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error while reading message: %v\n", err)
				return
			}

			var notification Notification
			err = json.Unmarshal(msg.Value, &notification)
			if err != nil {
				log.Printf("Failed to unmarshal JSON: %v\n", err)
				continue
			}

			userID := string(msg.Key)
			store.Add(userID, notification)

			err = ws.WriteJSON(notification)

			if err != nil {
				log.Printf("Failed to send message: %v\n", err)
				return
			}
		}
	}
}

func handleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "No notifications found for user",
			"notifications": []Notification{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notifications": notes,
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	store := &NotificationStore{
		data: make(UserNotification),
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, store)
	})

	router.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Error during upgrade: %v", err)
			return
		}
		log.Println("WebSocket connection established")
		defer ws.Close()

		userID := c.Query("userID")
		log.Printf("User ID: %s", userID)

		done := make(chan struct{})

		ctx := c.Request.Context()
		go readMessages(ctx, store, ws, done)

		<-ctx.Done()
		close(done)
		log.Println("WebSocket connection closed")
	})

	log.Println("WebSocket server started at ", ConsumerPort)
	if err := router.Run(ConsumerPort); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
