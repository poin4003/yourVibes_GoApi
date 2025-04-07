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

type PostEventPublisher struct {
	conn *rabbitmq.Connection
}

func NewPostEventPublisher(conn *rabbitmq.Connection) *PostEventPublisher {
	return &PostEventPublisher{
		conn: conn,
	}
}

func (p *PostEventPublisher) PublishNotification(ctx context.Context, msg interface{}, routingKey string) error {
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

func (p *PostEventPublisher) PublishStatistic(ctx context.Context, msg interface{}, routingKey string) error {
	body, err := json.Marshal(msg)
	if err != nil {
		global.Logger.Error("failed to marshal statistic payload", zap.Error(err))
		return err
	}

	err = p.conn.Publish(ctx, consts.StatisticsExName, routingKey, body, nil)
	if err != nil {
		global.Logger.Error("failed to publish statistic", zap.Error(err))
		return err
	}
	return nil
}
