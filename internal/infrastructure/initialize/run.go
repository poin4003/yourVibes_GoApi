package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
)

func Run() *gin.Engine {
	LoadConfig()
	InitLogger()
	routerGroup := InitDependencyInjection(
		InitPostgreSql(),
		InitRabbitMQ(),
		InitRedis(),
		socket_hub.NewNotificationSocketHub(),
		socket_hub.NewMessageSocketHub(),
	)
	response.InitCustomCode()

	r := InitRouter(*routerGroup)

	return r
}
