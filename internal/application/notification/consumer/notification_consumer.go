package consumer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"go.uber.org/zap"
	"strings"
)

type NotificationConsumer struct {
	queueName           string
	notificationService services.INotificationMQ
}

func NewNotificationConsumer(queueName string, service services.INotificationMQ) *NotificationConsumer {
	return &NotificationConsumer{queueName: queueName, notificationService: service}
}

func (c *NotificationConsumer) StartNotificationConsuming(ctx context.Context) error {
	ch, err := global.RabbitMQConn.GetChannel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	routingKeys := []string{
		"notification.bulk.db_websocket",
		"notification.single.db_websocket",
	}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			q.Name,
			key,
			consts.NotificationExName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				var notifMsg command.NotificationCommand
				if err = json.Unmarshal(msg.Body, &notifMsg); err != nil {
					global.Logger.Error("Failed to unmarshal notification command", zap.Error(err))
					break
				}

				parts := strings.Split(msg.RoutingKey, ".")
				if len(parts) < 3 || parts[0] != "notification" {
					global.Logger.Warn("Invalid routing key", zap.String("routing_key", msg.RoutingKey))
				}
				scope := parts[1]
				actions := strings.Split(parts[2], "_")

				switch scope {
				case "bulk":
					if err = c.notificationService.HandleBulkNotification(ctx, notifMsg, actions); err != nil {
						global.Logger.Error("Failed to handle bulk notification", zap.Error(err))
					}
				case "single":
					if err = c.notificationService.HandleSingleNotification(ctx, notifMsg, actions); err != nil {
						global.Logger.Error("Failed to handle single notification", zap.Error(err))
					}
				default:
					global.Logger.Warn("Unknown scope in routing key", zap.String("scope", scope))
				}
			case <-ctx.Done():
				global.Logger.Info("Stopping consumer")
				return
			}
		}
	}()

	global.Logger.Info("Consumer started", zap.String("queue", c.queueName))
	return nil
}

func InitNotificationConsumer(queueName string, service services.INotificationMQ) {
	consumer := NewNotificationConsumer(queueName, service)
	go func() {
		if err := consumer.StartNotificationConsuming(context.Background()); err != nil {
			global.Logger.Error("Failed to start notification consumer", zap.Error(err))
		}
	}()
	global.Logger.Info("Notification consumer started", zap.String("queue", queueName))
}
