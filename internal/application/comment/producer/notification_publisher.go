package producer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type NotificationPublisher struct {
	conn *rabbitmq.Connection
}

func NewNotificationPublisher(conn *rabbitmq.Connection) *NotificationPublisher {
	return &NotificationPublisher{
		conn: conn,
	}
}

func (p *NotificationPublisher) PublishNotification(ctx context.Context, msg interface{}, routingKey string) error {
	body, err := json.Marshal(msg)
	if err != nil {
		global.Logger.Error("failed to marshal notification payload", zap.Error(err))
		return err
	}

	if err = p.conn.Publish(ctx, consts.NotificationExName, routingKey, body,
		amqp091.Table{"original_routing_key": routingKey},
	); err != nil {
		global.Logger.Error("failed to publish notification", zap.Error(err))
		return err
	}

	return nil
}
