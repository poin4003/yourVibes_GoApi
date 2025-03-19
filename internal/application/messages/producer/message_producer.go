package producer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"go.uber.org/zap"
)

type MessagePublisher struct {
	conn *rabbitmq.Connection
}

func NewMessagePublisher(conn *rabbitmq.Connection) *MessagePublisher {
	return &MessagePublisher{
		conn: conn,
	}
}

func (p *MessagePublisher) PublishMessage(
	ctx context.Context,
	msg *command.CreateMessageCommand,
) error {
	body, err := json.Marshal(msg)
	if err != nil {
		global.Logger.Error("failed to marshal message", zap.Error(err))
	}

	err = p.conn.Publish(ctx, consts.MessageExName, "message.created", body)
	if err != nil {
		global.Logger.Error("failed to publish message", zap.Error(err))
		return err
	}

	return nil
}
