package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/settings"
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn        *amqp091.Connection
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

	for _, ex := range c.config.Exchanges {
		err = ch.ExchangeDeclare(ex.Name, ex.Type, true, false, false, false, nil)
		if err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to declare exchange %s: %v", ex.Name, err)
		}
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.conn.NotifyClose(c.notifyClose)
	c.mu.Unlock()

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
						log.Println("Reconnected to RabbitMQ successfully")
						return
					}
					attempts++
					log.Printf("Reconnect attempt %d/%d failed: %v", attempts, c.config.MaxReconnectAttempts, err)
					time.Sleep(2 * time.Second)
				}
				log.Fatalf("Failed to reconnect to RabbitMQ after %d attempts", c.config.MaxReconnectAttempts)
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
	if c.conn != nil {
		c.conn.Close()
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.channel == nil {
		return nil, fmt.Errorf("channel is nil")
	}
	return c.channel, nil
}
