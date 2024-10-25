package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

type UserNotifications map[string][]Notification

type NotificationStore struct {
	data UserNotifications
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

func readMessages(ctx context.Context, store *NotificationStore) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{KafkaServerAddress},
		Topic:   ConsumerTopic,
		GroupID: ConsumerGroup,
	})

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			return
		}

		var notification Notification
		err = json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			continue
		}

		userID := string(msg.Key)
		store.Add(userID, notification)
	}
}

func handleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
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

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}

func main() {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go readMessages(ctx, store)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("Failed to run the server: %v", err)
	}
}
