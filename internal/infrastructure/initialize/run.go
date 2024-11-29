package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	LoadConfig()
	InitCustomValidator()
	m := global.Config.PostgreSql
	fmt.Println("Loading configuration postgreSql", m.Username, m.Port)
	InitLogger()
	global.Logger.Info("Config log ok!!", zap.String("ok", "success"))
	InitCloudinary()
	InitRedis()
	InitPostgreSql()
	InitSocketHub()
	InitServiceInterface()
	InitCronJob()

	r := InitRouter()

	return r
}
