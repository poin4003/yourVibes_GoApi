package rabbitmq

import (
	"context"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"log"
	"net"
	"sync"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/settings"
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	Conn        *amqp091.Connection
	channel     *amqp091.Channel
	config      settings.RabbitMQSetting
	mu          sync.RWMutex
	notifyClose chan *amqp091.Error
}

func NewConnection(cfg settings.RabbitMQSetting) (*Connection, error) {
	c := &Connection{
		config:      cfg,
		notifyClose: make(chan *amqp091.Error),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	if err := c.setupExchanges(); err != nil {
		c.Close()
		return nil, fmt.Errorf("setup exchanges failed: %v", err)
	}

	go c.handleReconnect()
	return c, nil
}

func (c *Connection) connect() error {
	url := c.config.URL
	if url == "" {
		url = fmt.Sprintf("amqp://%s:%s@localhost:5672/%s", c.config.Username, c.config.Password, c.config.Vhost)
	}

	conn, err := amqp091.DialConfig(url, amqp091.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Duration(c.config.ConnectionTimeout)*time.Second)
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	c.mu.Lock()
	c.Conn = conn
	c.channel = ch
	c.Conn.NotifyClose(c.notifyClose)
	c.mu.Unlock()

	return nil
}

func (c *Connection) setupExchanges() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel == nil {
		return fmt.Errorf("channel is nil, cannot setup exchange")
	}

	err := c.channel.ExchangeDeclare(
		consts.NotificationExName,
		consts.NotificationExType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to setup exchange: %v", err)
	}

	return nil
}

func (c *Connection) handleReconnect() {
	for {
		select {
		case err := <-c.notifyClose:
			if err != nil {
				log.Printf("RabbitMQ connection closed: %v. Reconnecting...", err)
				attempts := 0
				for attempts < c.config.MaxReconnectAttempts {
					if err := c.connect(); err == nil {
						if err := c.setupExchanges(); err != nil {
							log.Printf("Failed to setup exchanges after reconnect: %v", err)
							continue
						}
						log.Println("Reconnected to RabbitMQ successfully")
						break
					}
					attempts++
					log.Printf("Reconnect attempt %d/%d failed: %v", attempts, c.config.MaxReconnectAttempts, err)
					time.Sleep(2 * time.Second)
				}
				if attempts >= c.config.MaxReconnectAttempts {
					log.Printf("Failed to reconnect to RabbitMQ after %d attempts", c.config.MaxReconnectAttempts)
					panic(fmt.Sprintf("Failed to reconnect to RabbitMQ after %d attempts", c.config.MaxReconnectAttempts))
				}
			}
		}
	}
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.channel != nil {
		c.channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *Connection) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.channel == nil {
		return fmt.Errorf("channel is nil, connection may be closed")
	}
	return c.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (c *Connection) GetChannel() (*amqp091.Channel, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil || c.Conn.IsClosed() || c.channel == nil {
		log.Println("RabbitMQ connection or channel closed, attempting to reconnect")
		if err := c.connect(); err != nil {
			return nil, fmt.Errorf("failed to reconnect to RabbitMQ: %v", err)
		}
		if err := c.setupExchanges(); err != nil {
			return nil, fmt.Errorf("failed to setup exchanges after reconnect: %v", err)
		}
	}

	return c.channel, nil
}
