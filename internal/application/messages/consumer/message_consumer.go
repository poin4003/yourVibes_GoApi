package consumer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MessageConsumer struct {
	messageService services.IMessageMQ
	conn           *rabbitmq.Connection
}

func NewMessageConsumer(
	service services.IMessageMQ,
	conn *rabbitmq.Connection,
) *MessageConsumer {
	return &MessageConsumer{
		messageService: service,
		conn:           conn,
	}
}

func (c *MessageConsumer) StartMessageConsuming(ctx context.Context) error {
	ch, err := c.conn.GetChannel()
	if err != nil {
		return err
	}

	// Declare DLX
	if err = ch.ExchangeDeclare(consts.MessageDLXName,
		"direct", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare message DLX", zap.Error(err))
		return err
	}

	// Declare exchange
	if err = ch.ExchangeDeclare(consts.MessageExName,
		"direct", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare message exchange", zap.Error(err))
		return err
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(consts.MessageDLQ, true, false, false, false,
		amqp091.Table{
			"x-message-ttl": int32(300000),
			"x-max-length":  int32(10000),
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare DLQ", zap.Error(err))
		return err
	}

	// Bind DLQ and DLX
	if err = ch.QueueBind(consts.MessageDLQ,
		"dlq_routing_key", consts.MessageDLXName, false, nil,
	); err != nil {
		global.Logger.Error("Failed to bind DLQ", zap.Error(err))
		return err
	}

	// Declare queue
	_, err = ch.QueueDeclare(consts.MessageQueue,
		true, false, false, false,
		amqp091.Table{
			"x-message-ttl":             int32(300000),
			"x-dead-letter-exchange":    consts.MessageDLXName,
			"x-dead-letter-routing-key": "dlq_routing_key",
			"x-max-length":              int32(10000),
			"x-overflow":                "reject-publish-dlx",
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare message queue", zap.Error(err))
		return err
	}

	// Bind Exchange and queue
	if err = ch.QueueBind(consts.MessageQueue,
		"message.created", consts.MessageExName, false, nil,
	); err != nil {
		global.Logger.Error("Failed to bind message queue", zap.Error(err))
		return err
	}

	msgsMain, err := ch.Consume(consts.MessageQueue, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume from main queue", zap.Error(err))
		return err
	}

	msgsDLQ, err := ch.Consume(consts.MessageDLQ, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume from DLQ", zap.Error(err))
		return err
	}

	global.Logger.Info("Message consumer started successfully", zap.String("queue", consts.MessageQueue))
	go c.consumeMessages(ctx, msgsMain, false)
	go c.consumeMessages(ctx, msgsDLQ, true)
	return nil
}

func (c *MessageConsumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery, isDLQ bool) {
	queueName := consts.MessageQueue
	if isDLQ {
		queueName = consts.MessageDLQ
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
				var msgCommand command.CreateMessageCommand
				if err := json.Unmarshal(msg.Body, &msgCommand); err != nil {
					global.Logger.Error("Failed to unmarshal command", zap.Error(err))
					msg.Ack(false)
					continue
				}

				routingKey := msg.RoutingKey
				if headerKey, ok := msg.Headers["original_routing_key"]; ok {
					if key, ok := headerKey.(string); ok {
						routingKey = key
					}
				}
				if routingKey != "message.created" {
					global.Logger.Warn("Invalid routing key format", zap.String("routing_key", routingKey))
				}

				if err := c.messageService.HandleMessage(ctx, &msgCommand); err != nil {
					global.Logger.Error("Failed to handle message", zap.Error(err))
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
			}
		case <-ctx.Done():
			global.Logger.Info("Consumer is shutting down", zap.String("queue", queueName))
			return
		}
	}
}

func (c *MessageConsumer) processDLQMessage(ctx context.Context, msg amqp091.Delivery) {
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

	global.Logger.Info("Processing DLQ message", zap.Int("retry_count", count), zap.String("queue", consts.MessageQueue), zap.String("routing_key", msg.RoutingKey))

	if count < 3 {
		var msgCommand command.CreateMessageCommand
		if err := json.Unmarshal(msg.Body, &msgCommand); err != nil {
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
		if routingKey != "message.created" {
			global.Logger.Warn("Invalid routing key format in DLQ", zap.String("routing_key", routingKey))
		}

		if err := c.messageService.HandleMessage(ctx, &msgCommand); err != nil {
			global.Logger.Error("Failed to handle DLQ message", zap.Error(err))
			msg.Nack(false, true)
			return
		}

		err := c.republishMessage(msg, consts.MessageQueue)
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

func (c *MessageConsumer) republishMessage(msg amqp091.Delivery, queue string) error {
	ch, err := c.conn.GetChannel()
	if err != nil {
		return err
	}

	routingKey := "message.created"
	if headerKey, ok := msg.Headers["original_routing_key"]; ok {
		if key, ok := headerKey.(string); ok {
			routingKey = key
		}
	}

	err = ch.Publish(consts.MessageExName,
		routingKey, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			Headers:     msg.Headers,
		},
	)
	return err
}

func InitMessageConsumer(messageService services.IMessageMQ, conn *rabbitmq.Connection) {
	consumer := NewMessageConsumer(messageService, conn)
	if err := consumer.StartMessageConsuming(context.Background()); err != nil {
		global.Logger.Error("Failed to start message consuming", zap.Error(err))
	} else {
		global.Logger.Info("Message consumer initialized successfully", zap.String("queue", consts.MessageQueue))
	}
}
