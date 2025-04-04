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
	Conn       *amqp091.Connection
	channel    *amqp091.Channel
	config     settings.RabbitMQSetting
	mu         sync.RWMutex
	notifyConn chan *amqp091.Error
	notifyChan chan *amqp091.Error
}

func NewConnection(cfg settings.RabbitMQSetting) (*Connection, error) {
	c := &Connection{
		config:     cfg,
		notifyConn: make(chan *amqp091.Error),
		notifyChan: make(chan *amqp091.Error),
	}

	if err := c.connectWithRetry(); err != nil {
		return nil, err
	}

	go c.handleReconnect()
	return c, nil
}

func (c *Connection) connectWithRetry() error {
	attempts := 0
	maxAttempts := c.config.MaxReconnectAttempts
	url := c.config.URL
	if url == "" {
		url = fmt.Sprintf("amqp://%s:%s@localhost:5672/%s", c.config.Username, c.config.Password, c.config.Vhost)
	}

	for attempts < maxAttempts {
		conn, err := amqp091.DialConfig(url, amqp091.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(c.config.ConnectionTimeout)*time.Second)
			},
		})
		if err == nil {
			ch, err := c.openChannelWithRetry(conn)
			if err == nil {
				c.mu.Lock()
				if c.Conn != nil {
					c.Conn.Close()
				}
				c.Conn = conn
				c.channel = ch
				c.Conn.NotifyClose(c.notifyConn)
				c.channel.NotifyClose(c.notifyChan)
				c.mu.Unlock()
				log.Println("Connected to RabbitMQ successfully")
				return nil
			}
			conn.Close()
		}
		attempts++
		log.Printf("Connection attempt %d/%d failed: %v", attempts, maxAttempts, err)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to connect to RabbitMQ after %d attempts", maxAttempts)
}

func (c *Connection) openChannelWithRetry(conn *amqp091.Connection) (*amqp091.Channel, error) {
	attempts := 0
	maxAttempts := 5

	for attempts < maxAttempts {
		ch, err := conn.Channel()
		if err == nil {
			return ch, nil
		}
		attempts++
		log.Printf("Channel open attempt %d/%d failed: %v", attempts, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to open channel after %d attempts", maxAttempts)
}

func (c *Connection) handleReconnect() {
	for {
		select {
		case err := <-c.notifyConn:
			if err != nil {
				log.Printf("RabbitMQ connection closed: %v. Reconnecting...", err)
				c.mu.Lock()
				if c.Conn != nil {
					c.Conn.Close()
				}
				c.Conn = nil
				c.channel = nil
				c.mu.Unlock()
				if err := c.connectWithRetry(); err != nil {
					log.Printf("Reconnect failed: %v", err)
				}
			}
		case err := <-c.notifyChan:
			if err != nil {
				log.Printf("RabbitMQ channel closed: %v. Reconnecting...", err)
				c.mu.Lock()
				c.channel = nil
				c.mu.Unlock()
				if err := c.connectWithRetry(); err != nil {
					log.Printf("Reconnect failed: %v", err)
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

func (c *Connection) Publish(ctx context.Context, exchange, routingKey string, body []byte, headers amqp091.Table) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.channel == nil {
		return fmt.Errorf("channel is nil, connection may be closed")
	}
	return c.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
	})
}

func (c *Connection) GetChannel() (*amqp091.Channel, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil || c.Conn.IsClosed() || c.channel == nil {
		log.Println("RabbitMQ connection or channel closed, attempting to reconnect")
		if err := c.connectWithRetry(); err != nil {
			return nil, fmt.Errorf("failed to reconnect to RabbitMQ: %v", err)
		}
	}

	return c.channel, nil
}
