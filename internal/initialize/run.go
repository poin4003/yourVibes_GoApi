package initialize

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
)

func Run() {
	LoadConfig()
	m := global.Config.PostgreSql
	fmt.Println("Loading configuration postgreSql", m.Username, m.Port)
	InitLogger()
	global.Logger.Info("Config log ok!!", zap.String("ok", "success"))
	InitRedis()
	InitPostgreSqlc()

	r := InitRouter()

	r.Run(":8080")
}
