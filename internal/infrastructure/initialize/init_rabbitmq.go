package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"go.uber.org/zap"
)

func InitRabbitMQ() {
	rabbitMQConn, err := rabbitmq.NewConnection(global.Config.RabbitMQSetting)
	if err != nil {
		global.Logger.Error("init rabbitmq connection error", zap.Error(err))
	}
	global.Logger.Info("init rabbitmq connection success")
	global.RabbitMQConn = rabbitMQConn
}
