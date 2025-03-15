package initialize

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	LoadConfig()
	m := global.Config.PostgreSql
	fmt.Println("Loading configuration postgreSql", m.Username, m.Port)
	InitLogger()
	global.Logger.Info("Config log ok!!", zap.String("ok", "success"))
	InitRedis()
	InitRabbitMQ()
	defer func() {
		if global.RabbitMQConn != nil {
			global.RabbitMQConn.Close()
			global.Logger.Info("RabbitMQ connection closed")
		}
	}()

	InitPostgreSql()
	InitSocketHub()
	InitServiceInterface()
	InitCronJob()
	response.InitCustomCode()

	r := InitRouter()

	return r
}
