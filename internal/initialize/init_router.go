package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/routers"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine

	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.New()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// middlewares
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},                                                // Cho phép tất cả các origin
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Cho phép các phương thức HTTP
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},          // Cho phép các header cụ thể
		ExposeHeaders: []string{"Content-Length"},
	}))
	//r.Use() // logging
	//r.Use() // cross
	//r.Use() // limiter global
	adminRouter := routers.RouterGroupApp.Admin
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/v1/2024")
	{
		MainGroup.GET("/checkStatus", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "Ok",
			})
		})
	}
	{
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitPostRouter(MainGroup)
	}
	{
		adminRouter.InitAdminRouter(MainGroup)
	}
	{
	}
	return r
}
