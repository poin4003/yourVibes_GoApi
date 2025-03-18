package consumer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"strings"
)

type NotificationConsumer struct {
	queueName           string
	dlqName             string
	notificationService services.INotificationMQ
}

func NewNotificationConsumer(
	queueName string,
	dlqName string,
	service services.INotificationMQ,
) *NotificationConsumer {
	return &NotificationConsumer{
		queueName:           queueName,
		dlqName:             dlqName,
		notificationService: service,
	}
}

func (c *NotificationConsumer) StartNotificationConsuming(ctx context.Context) error {
	ch, err := global.RabbitMQConn.GetChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(c.dlqName, true, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to declare DLQ", zap.Error(err))
		return err
	}

	err = ch.QueueBind(
		c.dlqName,
		"dlq_routing_key",
		consts.NotificationDLXName,
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to bind DLQ", zap.Error(err))
		return err
	}

	_, err = ch.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    consts.NotificationDLXName,
			"x-dead-letter-routing-key": "dlq_routing_key",
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare notification queue", zap.Error(err))
		return err
	}

	routingKeys := []string{
		"notification.bulk.db_websocket",
		"notification.single.db_websocket",
	}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			c.queueName,
			key,
			consts.NotificationExName,
			false,
			nil,
		)
		if err != nil {
			global.Logger.Error("Failed to bind notification queue", zap.Error(err), zap.String("routing_key", key))
			return err
		}
	}

	qMain, err := ch.QueueDeclarePassive(c.queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgsMain, err := ch.Consume(qMain.Name, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume from main queue", zap.Error(err))
		return err
	}

	qDLQ, err := ch.QueueDeclarePassive(c.dlqName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgsDLQ, err := ch.Consume(qDLQ.Name, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume from DLQ", zap.Error(err))
		return err
	}

	global.Logger.Info("Notification consumer started successfully", zap.String("queue", c.queueName))
	go c.consumeMessages(ctx, msgsMain, false)
	go c.consumeMessages(ctx, msgsDLQ, true)
	return nil
}

func (c *NotificationConsumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery, isDLQ bool) {
	queueName := c.queueName
	if isDLQ {
		queueName = c.dlqName
	}
	global.Logger.Info("Consumer goroutine started", zap.String("queue", queueName), zap.Bool("isDLQ", isDLQ))
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				global.Logger.Warn("Message channel closed", zap.String("queue", queueName))
				return
			}
			global.Logger.Info("Received message", zap.String("queue", queueName), zap.String("routing_key", msg.RoutingKey))

			if isDLQ {
				c.processDLQMessage(ctx, msg)
			} else {
				var notifMsg command.NotificationCommand
				if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
					global.Logger.Error("Failed to unmarshal notification command", zap.Error(err))
					msg.Nack(false, true)
					continue
				}

				routingKey := msg.RoutingKey
				if headerKey, ok := msg.Headers["original_routing_key"]; ok {
					if key, ok := headerKey.(string); ok {
						routingKey = key
					}
				}
				parts := strings.Split(routingKey, ".")
				if len(parts) < 3 || parts[0] != "notification" {
					global.Logger.Warn("Invalid routing key format", zap.String("routing_key", routingKey))
				}
				scope := ""
				actions := []string{}
				if len(parts) > 1 {
					scope = parts[1]
					if len(parts) > 2 {
						actions = strings.Split(parts[2], "_")
					}
				}

				switch scope {
				case "bulk":
					if err := c.notificationService.HandleBulkNotification(ctx, notifMsg, actions); err != nil {
						global.Logger.Error("Failed to handle bulk notification", zap.Error(err))
						msg.Nack(false, true)
						continue
					}
				case "single":
					if err := c.notificationService.HandleSingleNotification(ctx, notifMsg, actions); err != nil {
						global.Logger.Error("Failed to handle single notification", zap.Error(err))
						msg.Nack(false, true)
						continue
					}
				default:
					global.Logger.Warn("Unknown scope in routing key", zap.String("scope", scope))
				}

				msg.Ack(false)
			}
		case <-ctx.Done():
			global.Logger.Info("Consumer is shutting down", zap.String("queue", queueName))
			return
		}
	}
}

func (c *NotificationConsumer) processDLQMessage(ctx context.Context, msg amqp091.Delivery) {
	count := 0
	if headers, ok := msg.Headers["x-death"]; ok {
		if deaths, ok := headers.([]interface{}); ok && len(deaths) > 0 {
			if death, ok := deaths[0].(amqp091.Table); ok {
				if c, ok := death["count"]; ok {
					if countInt, ok := c.(int32); ok {
						count = int(countInt)
					}
				}
			}
		}
	}

	global.Logger.Info("Processing DLQ message", zap.Int("retry_count", count), zap.String("queue", c.dlqName), zap.String("routing_key", msg.RoutingKey))

	if count < 3 {
		var notifMsg command.NotificationCommand
		if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
			global.Logger.Error("Failed to unmarshal DLQ command", zap.Error(err))
			msg.Nack(false, true)
			return
		}

		routingKey := msg.RoutingKey
		if headerKey, ok := msg.Headers["original_routing_key"]; ok {
			if key, ok := headerKey.(string); ok {
				routingKey = key
			}
		}
		parts := strings.Split(routingKey, ".")
		if len(parts) < 3 || parts[0] != "notification" {
			global.Logger.Warn("Invalid routing key format in DLQ", zap.String("routing_key", routingKey))
		}
		scope := ""
		actions := []string{}
		if len(parts) > 1 {
			scope = parts[1]
			if len(parts) > 2 {
				actions = strings.Split(parts[2], "_")
			}
		}

		var err error
		switch scope {
		case "bulk":
			err = c.notificationService.HandleBulkNotification(ctx, notifMsg, actions)
		case "single":
			err = c.notificationService.HandleSingleNotification(ctx, notifMsg, actions)
		default:
			global.Logger.Warn("Unknown scope in routing key", zap.String("scope", scope))
		}

		if err != nil {
			global.Logger.Error("Failed to handle DLQ notification", zap.Error(err))
			msg.Nack(false, true)
			return
		}

		err = c.republishMessage(msg, c.queueName)
		if err != nil {
			global.Logger.Error("Failed to republish message to main queue", zap.Error(err))
			msg.Nack(false, true)
			return
		}
		msg.Ack(false)
	} else {
		global.Logger.Warn("Max retry reached, discarding message", zap.String("message", string(msg.Body)))
		msg.Ack(false)
	}
}

func (c *NotificationConsumer) republishMessage(msg amqp091.Delivery, queue string) error {
	ch, err := global.RabbitMQConn.GetChannel()
	if err != nil {
		return err
	}

	routingKey := msg.RoutingKey
	if headerKey, ok := msg.Headers["original_routing_key"]; ok {
		if key, ok := headerKey.(string); ok {
			routingKey = key
		}
	}

	err = ch.Publish(
		consts.NotificationExName,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			Headers:     msg.Headers,
		},
	)
	return err
}

func InitNotificationConsumer(queueName string, dlq string, service services.INotificationMQ) {
	consumer := NewNotificationConsumer(queueName, dlq, service)
	go func() {
		if err := consumer.StartNotificationConsuming(context.Background()); err != nil {
			global.Logger.Error("Failed to start notification consumer", zap.Error(err))
		} else {
			global.Logger.Info("Notification consumer initialized successfully", zap.String("queue", queueName))
		}
	}()
}
