package consumer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/statistic/services"
	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type StatisticConsumer struct {
	queueName        string
	exchangeName     string
	batch            map[uuid.UUID]*command.UpsertStatisticCommand
	mu               sync.Mutex
	cron             *cron.Cron
	statisticService services.IStatisticMQ
}

func NewStatisticConsumer(
	queueName string,
	exchangeName string,
	statisticService services.IStatisticMQ,
) *StatisticConsumer {
	c := &StatisticConsumer{
		queueName:        queueName,
		exchangeName:     exchangeName,
		batch:            make(map[uuid.UUID]*command.UpsertStatisticCommand),
		cron:             cron.New(),
		statisticService: statisticService,
	}

	_, err := c.cron.AddFunc("@every 5m", func() {
		c.processBatch(context.Background())
	})
	if err != nil {
		global.Logger.Error("", zap.Error(err))
	}

	return c
}

func (c *StatisticConsumer) StartStatisticsConsuming(ctx context.Context) error {
	ch, err := global.RabbitMQConn.GetChannel()
	if err != nil {
		return err
	}

	// Declare Statistic queue
	_, err = ch.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("Failed to declare a queue", zap.Error(err))
		return err
	}
	err = ch.QueueBind(c.queueName, "stats.post", c.exchangeName, false, nil)
	if err != nil {
		global.Logger.Error("Failed to bind a queue", zap.Error(err))
		return err
	}

	// Consume
	msgs, err := ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		global.Logger.Error("Failed to consume a statistic queue", zap.Error(err))
		return err
	}

	c.cron.Start()
	global.Logger.Info("Started cron scheduler for batch processing")
	global.Logger.Info("Consuming statistic queue", zap.String("queue", c.queueName))
	go c.consumeMessages(ctx, msgs)
	return nil
}

func (c *StatisticConsumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery) {
	global.Logger.Info("Consume statistic started", zap.String("queue", c.queueName))

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				global.Logger.Info("Message channel closed", zap.String("queue", c.queueName))
				return
			}

			var eventCommand command.EventCommand
			if err := json.Unmarshal(msg.Body, &eventCommand); err != nil {
				global.Logger.Error("Failed to unmarshal event", zap.Error(err))
				msg.Ack(false)
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
		case <-ctx.Done():
			global.Logger.Info("Consume statistic stopped", zap.String("queue", c.queueName))
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
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
	global.Logger.Info("Finished processing batch periodically")
}

func InitStatisticsConsumer(
	queueName string,
	exchangeName string,
	statisticService services.IStatisticMQ,
) {
	consumer := NewStatisticConsumer(
		queueName,
		exchangeName,
		statisticService,
	)
	go func() {
		if err := consumer.StartStatisticsConsuming(context.Background()); err != nil {
			global.Logger.Error("Failed to start statistics consuming", zap.Error(err))
		} else {
			global.Logger.Info("Statistics consumer initialized successfully", zap.String("queue", queueName))
		}
	}()
}
