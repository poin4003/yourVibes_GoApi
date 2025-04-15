package consumer

import (
	"context"
	"encoding/json"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type PostModerationConsumer struct {
	postService services.IPostUser
	conn        *rabbitmq.Connection
}

func NewPostModerationConsumer(
	postService services.IPostUser,
	conn *rabbitmq.Connection,
) *PostModerationConsumer {
	return &PostModerationConsumer{
		postService: postService,
		conn:        conn,
	}
}

func (c *PostModerationConsumer) StartPostModerationConsuming(ctx context.Context) error {
	ch, err := c.conn.GetChannel()
	if err != nil {
		return err
	}

	// Declare DLX
	if err = ch.ExchangeDeclare(consts.CreatePostDLX,
		"direct", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare create_post_dlx", zap.Error(err))
		return err
	}

	// Declare exchange
	if err = ch.ExchangeDeclare(consts.CreatePostExchange,
		"direct", true, false, false, false, nil,
	); err != nil {
		global.Logger.Error("failed to declare create_post_exchange", zap.Error(err))
		return err
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(consts.CreatePostDLQ,
		true, false, false, false,
		amqp091.Table{
			"x-message-ttl": int32(3600000),
			"x-max-length":  int32(10000),
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare create_post_dlq", zap.Error(err))
		return err
	}

	if err = ch.QueueBind(
		consts.CreatePostDLQ,
		consts.CreatePostDLQ,
		consts.CreatePostDLX,
		false,
		nil,
	); err != nil {
		global.Logger.Error("Failed to bind create_post_dlq", zap.Error(err))
		return err
	}

	// Declare main queue
	_, err = ch.QueueDeclare(
		consts.CreatePostQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(3600000),
			"x-max-length":              int32(10000),
			"x-dead-letter-exchange":    consts.CreatePostDLX,
			"x-dead-letter-routing-key": consts.CreatePostDLQ,
			"x-overflow":                "reject-publish-dlx",
		},
	)
	if err != nil {
		global.Logger.Error("Failed to declare create_post_queue", zap.Error(err))
		return err
	}

	if err = ch.QueueBind(
		consts.CreatePostQueue,
		consts.CreatePostQueue,
		consts.CreatePostExchange,
		false,
		nil,
	); err != nil {
		global.Logger.Error("Failed to bind create_post_queue", zap.Error(err))
		return err
	}

	// Consume messages from main queue
	msgsMain, err := ch.Consume(consts.CreatePostQueue,
		"", false, false, false, false, nil,
	)
	if err != nil {
		global.Logger.Error("Failed to consume messages from create_post_queue", zap.Error(err))
		return err
	}

	// Consume messages from DLQ
	msgsDLQ, err := ch.Consume(consts.CreatePostDLQ,
		"", false, false, false, false, nil,
	)
	if err != nil {
		global.Logger.Error("Failed to consume messages from create_post_dlq", zap.Error(err))
		return err
	}

	global.Logger.Info("Post moderation consumer started successfully", zap.String("queue", consts.CreatePostQueue))
	go c.consumeMessages(ctx, msgsMain, false)
	go c.consumeMessages(ctx, msgsDLQ, true)
	return nil
}

func (c *PostModerationConsumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery, isDLQ bool) {
	queueName := consts.CreatePostQueue
	if isDLQ {
		queueName = consts.CreatePostDLQ
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
				var response command.PostModerationResult
				if err := json.Unmarshal(msg.Body, &response); err != nil {
					global.Logger.Error("Failed to unmarshal post moderation response", zap.Error(err))
					msg.Nack(false, true)
					continue
				}

				var err error
				if response.Media.Label != "normal" {
					rejectCommand := &command.RejectPostCommand{
						PostId: response.PostID,
						Label:  response.Media.Label,
					}
					err = c.postService.RejectPost(ctx, rejectCommand)
					if err != nil {
						global.Logger.Error("Failed to process post moderation", zap.Error(err))
						msg.Ack(false)
						continue
					}
				} else {
					approveCommand := &command.ApprovePostCommand{
						PostId:       response.PostID,
						CensoredText: response.Content.CensoredText,
					}
					err = c.postService.ApprovePost(ctx, approveCommand)
					if err != nil {
						global.Logger.Error("Failed to process post moderation", zap.Error(err))
						msg.Ack(false)
						continue
					}
				}

				msg.Ack(false)
			}

		case <-ctx.Done():
			global.Logger.Info("Consumer is shutting down", zap.String("queue", queueName))
			return
		}
	}
}

func (c *PostModerationConsumer) processDLQMessage(ctx context.Context, msg amqp091.Delivery) {
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

	global.Logger.Info("Processing DLQ message", zap.Int("retry_count", count), zap.String("queue", consts.CreatePostQueue), zap.String("routing_key", msg.RoutingKey))

	if count < 3 {
		var response command.PostModerationResult
		if err := json.Unmarshal(msg.Body, &response); err != nil {
			global.Logger.Error("Failed to unmarshal DLQ response", zap.Error(err))
			msg.Nack(false, true)
			return
		}

		var err error
		if response.Media.Label != "normal" {
			rejectCommand := &command.RejectPostCommand{
				PostId: response.PostID,
				Label:  response.Media.Label,
			}
			err = c.postService.RejectPost(ctx, rejectCommand)
			if err != nil {
				global.Logger.Error("Failed to process post moderation", zap.Error(err))
				msg.Nack(false, true)
			}
		} else {
			approveCommand := &command.ApprovePostCommand{
				PostId:       response.PostID,
				CensoredText: response.Content.CensoredText,
			}
			err = c.postService.ApprovePost(ctx, approveCommand)
			if err != nil {
				global.Logger.Error("Failed to process post moderation", zap.Error(err))
				msg.Nack(false, true)
			}
		}

		// Republish message to main queue
		if err = c.republishMessage(msg, consts.CreatePostQueue); err != nil {
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

func (c *PostModerationConsumer) republishMessage(msg amqp091.Delivery, queue string) error {
	ch, err := c.conn.GetChannel()
	if err != nil {
		return err
	}

	err = ch.Publish(
		consts.CreatePostExchange,
		consts.CreatePostQueue,
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

func InitPostModerationConsumer(postService services.IPostUser, conn *rabbitmq.Connection) {
	consumer := NewPostModerationConsumer(postService, conn)
	go func() {
		if err := consumer.StartPostModerationConsuming(context.Background()); err != nil {
			global.Logger.Error("Failed to start post moderation consumer", zap.Error(err))
		} else {
			global.Logger.Info("Post moderation consumer initialized successfully", zap.String("queue", consts.CreatePostQueue))
		}
	}()
}
