package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"

	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	LoadConfig()
	InitLogger()
	rdb := InitRedis()
	rabbitMQConnection := InitRabbitMQ()
	db := InitPostgreSql()
	global.MessageSocketHub = socket_hub.NewMessageSocketHub()
	global.NotificationSocketHub = socket_hub.NewNotificationSocketHub()
	InitDependencyInjection(db, rabbitMQConnection, rdb, global.NotificationSocketHub, global.MessageSocketHub)
	response.InitCustomCode()

	r := InitRouter()

	return r
}
