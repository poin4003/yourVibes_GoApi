package consumer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"go.uber.org/zap"
)

type MessageConsumer struct {
	messageService services.IMessageMQ
	queueName      string
}

func NewMessageConsumer(service services.IMessageMQ, queueName string) *MessageConsumer {
	return &MessageConsumer{
		messageService: service,
		queueName:      queueName,
	}
}

func (c *MessageConsumer) StartMessageConsuming(ctx context.Context) error {
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
		"",
		consts.MessageExName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				var msgCommand command.CreateMessageCommand
				if err = json.Unmarshal(msg.Body, &msgCommand); err != nil {
					global.Logger.Error("Failed to unmarshal command", zap.Error(err))
					break
				}

				if err = c.messageService.HandleMessage(ctx, &msgCommand); err != nil {
					global.Logger.Error("Failed to handle message", zap.Error(err))
				}
			case <-ctx.Done():
				global.Logger.Info("Message consumer is shutting down")
				return
			}
		}
	}()

	global.Logger.Info("Message consumer started", zap.String("queue", c.queueName))
	return nil
}

func InitMessageConsumer(queueName string, messageService services.IMessageMQ) {
	consumer := NewMessageConsumer(messageService, queueName)
	go func() {
		if err := consumer.StartMessageConsuming(context.Background()); err != nil {
			global.Logger.Error("Failed to start message consuming", zap.Error(err))
		}
	}()
	global.Logger.Info("Message consumer started", zap.String("queue", queueName))
}
