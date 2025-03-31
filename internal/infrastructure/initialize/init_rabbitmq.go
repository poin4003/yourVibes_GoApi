package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"go.uber.org/zap"
)

func InitRabbitMQ() *rabbitmq.Connection {
	rabbitMQConn, err := rabbitmq.NewConnection(global.Config.RabbitMQSetting)
	if err != nil {
		global.Logger.Error("init rabbitmq connection error", zap.Error(err))
		panic("Failed to connect to rabbitmq: " + err.Error())
	}
	global.Logger.Info("init rabbitmq connection success")

	if rabbitMQConn.Conn.IsClosed() {
		global.Logger.Error("rabbitmq connection is closed")
		panic("RabbitMQ connection closed immediately after initialization")
	}

	return rabbitMQConn
}
