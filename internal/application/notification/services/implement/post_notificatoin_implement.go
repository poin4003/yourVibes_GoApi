package implement

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"go.uber.org/zap"
)

type Message struct {
	Type    consts.NotificationType `json:"type"`
	PostID  string                  `json:"post_id,omitempty"`
	UserID  string                  `json:"user_id,omitempty"`
	Content string                  `json:"content,omitempty"`
}

type Consumer struct {
	queueName string
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	ch, err := global.RabbitMQConn.GetChannel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"notification.*",
		"notification_exchange",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				var notifMsg Message
				if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
					global.Logger.Error("Failed to unmarshal notification message", zap.Error(err))
					continue
				}
				c.handleMessage(notifMsg)
			case <-ctx.Done():
				global.Logger.Info("Stopping consumer")
				return
			}
		}
	}()

	global.Logger.Info("Consumer started", zap.String("queue", c.queueName))
	return nil
}

func (c *Consumer) handleMessage(msg Message) {
	global.Logger.Info("Received notification",
		zap.String("type", string(msg.Type)),
		zap.String("use_id", msg.UserID),
		zap.String("post_id", msg.PostID),
		zap.String("content", msg.Content),
	)

	if global.SocketHub != nil {
		socketMsg := struct {
			Type    string `json:"type"`
			PostID  string `json:"post_id,omitempty"`
			UserID  string `json:"user_id"`
			Content string `json:"content,omitempty"`
		}{
			Type:    string(msg.Type),
			PostID:  msg.PostID,
			UserID:  msg.UserID,
			Content: msg.Content,
		}
		if err := global.SocketHub.SendNotification(msg.UserID, socketMsg); err != nil {
			global.Logger.Error("Failed to send notification", zap.Error(err))
		}
	}
}
