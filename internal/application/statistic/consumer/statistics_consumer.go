package consumer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
)

type StatisticConsumer struct {
	batch            map[uuid.UUID]*command.UpsertStatisticCommand
	mu               sync.Mutex
	cron             *cron.Cron
	statisticService services.IStatisticMQ
	conn             *rabbitmq.Connection
}

func NewStatisticConsumer(
	statisticService services.IStatisticMQ,
	conn *rabbitmq.Connection,
) *StatisticConsumer {
	c := &StatisticConsumer{
		batch:            make(map[uuid.UUID]*command.UpsertStatisticCommand),
		cron:             cron.New(),
		statisticService: statisticService,
		conn:             conn,
	}

	_, err := c.cron.AddFunc("@every 2m", func() {
		c.processBatch(context.Background())
	})
	if err != nil {
		global.Logger.Error("", zap.Error(err))
	}

	return c
}

func (c *StatisticConsumer) StartStatisticsConsuming(ctx context.Context) error {
	ch, err := c.conn.GetChannel()
	if err != nil {
		return err
	}

	// Declare exchange
	if err = ch.ExchangeDeclare(consts.StatisticsExName,
		"topic", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare message statistics exchange", zap.Error(err))
		return err
	}

	// Declare DLX
	if err = ch.ExchangeDeclare(consts.StatisticDLXName,
		"topic", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare message statistics dlx", zap.Error(err))
		return err
	}

	// Declare dlq
	_, err = ch.QueueDeclare(consts.StatisticsDLQ,
		true, false, false, false,
		amqp091.Table{
			"x-message-ttl": int32(300000),
			"x-max-length":  int32(10000),
		},
	)
	if err != nil {
		global.Logger.Error("failed to declare message statistics queue", zap.Error(err))
		return err
	}

	// Declare Statistic queue
	_, err = ch.QueueDeclare(consts.StatisticsQueue,
		true, false, false, false,
		amqp091.Table{
			"x-message-ttl":             int32(300000),
			"x-dead-letter-exchange":    consts.StatisticDLXName,
			"x-dead-letter-routing-key": "stats.post",
			"x-max-length":              int32(10000),
			"x-overflow":                "reject-publish-dlx",
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare a queue", zap.Error(err))
		return err
	}

	// Bind dlq and dlx
	err = ch.QueueBind(consts.StatisticsDLQ, "stats.post", consts.StatisticDLXName, false, nil)
	if err != nil {
		global.Logger.Error("Failed to bind a queue", zap.Error(err))
		return err
	}
	// Bind queue and exchange
	err = ch.QueueBind(consts.StatisticsQueue, "stats.post", consts.StatisticsExName, false, nil)
	if err != nil {
		global.Logger.Error("Failed to bind a queue", zap.Error(err))
		return err
	}

	// Consume
	msgsMain, err := ch.Consume(consts.StatisticsQueue, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume a statistic queue", zap.Error(err))
		return err
	}

	msgsDlq, err := ch.Consume(consts.StatisticsDLQ, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume a DLQ queue", zap.Error(err))
		return err
	}

	c.cron.Start()
	global.Logger.Info("Started cron scheduler for batch processing")
	global.Logger.Info("Consuming statistic queue", zap.String("queue", consts.StatisticsQueue))
	go c.consumeMessages(ctx, msgsMain, false)
	go c.consumeMessages(ctx, msgsDlq, true)
	return nil
}

func (c *StatisticConsumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery, isDLQ bool) {
	global.Logger.Info("Consume statistic started", zap.String("queue", consts.StatisticsQueue))
	queueName := consts.MessageQueue
	if isDLQ {
		queueName = consts.MessageDLQ
	}

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				global.Logger.Info("Message channel closed", zap.String("queue", queueName))
				return
			}

			if isDLQ {
				c.processDLQMessage(ctx, msg)
			} else {
				var eventCommand command.EventCommand
				if err := json.Unmarshal(msg.Body, &eventCommand); err != nil {
					global.Logger.Error("Failed to unmarshal event", zap.Error(err))
					msg.Nack(false, false)
					continue
				}

				c.mu.Lock()
				if _, exists := c.batch[eventCommand.PostId]; !exists {
					c.batch[eventCommand.PostId] = &command.UpsertStatisticCommand{
						PostId: eventCommand.PostId,
					}
				}
				cmd := c.batch[eventCommand.PostId]
				switch eventCommand.EventType {
				case "reach":
					cmd.Reach += eventCommand.Count
				case "clicks":
					cmd.Clicks += eventCommand.Count
				case "impression":
					cmd.Impression += eventCommand.Count
				default:
					global.Logger.Warn("Unknown event type", zap.String("event_type", eventCommand.EventType))
				}
				c.mu.Unlock()

				msg.Ack(false)
			}
		case <-ctx.Done():
			global.Logger.Info("Consume statistic stopped", zap.String("queue", consts.StatisticsQueue))
			c.processBatch(ctx)
			c.cron.Stop()
			return
		}
	}
}

func (c *StatisticConsumer) processBatch(ctx context.Context) {
	c.mu.Lock()
	if len(c.batch) == 0 {
		c.mu.Unlock()
		return
	}

	batch := c.batch
	c.batch = make(map[uuid.UUID]*command.UpsertStatisticCommand)
	c.mu.Unlock()

	global.Logger.Info("Processing batch", zap.Int("batch_size", len(batch)))
	var wg sync.WaitGroup
	for postId, cmd := range batch {
		wg.Add(1)
		go func(pid uuid.UUID, cmd *command.UpsertStatisticCommand) {
			defer wg.Done()
			err := c.statisticService.UpsertStatistic(ctx, pid, cmd)
			if err != nil {
				global.Logger.Error("Failed to upsert statistic", zap.Error(err))
			}
		}(postId, cmd)
	}
	wg.Wait()
	global.Logger.Info("Finished processing batch periodically")
}

func (c *StatisticConsumer) processDLQMessage(ctx context.Context, msg amqp091.Delivery) {
	global.Logger.Warn("Message moved to DLQ", zap.ByteString("body", msg.Body))
	msg.Ack(false)
}

func InitStatisticsConsumer(
	statisticService services.IStatisticMQ,
	conn *rabbitmq.Connection,
) {
	consumer := NewStatisticConsumer(
		statisticService,
		conn,
	)
	if err := consumer.StartStatisticsConsuming(context.Background()); err != nil {
		global.Logger.Error("Failed to start statistics consuming", zap.Error(err))
	} else {
		global.Logger.Info("Statistics consumer initialized successfully", zap.String("queue", consts.StatisticsQueue))
	}
}
