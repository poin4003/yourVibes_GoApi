package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers"
)

func InitRouter(routerGroup routers.RouterGroup) *gin.Engine {
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
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	//r.Use() // logging
	//r.Use() // limiter global

	r.Use(middlewares.ErrorHandlerMiddleware())

	adminRouter := routerGroup.Admin
	userRouter := routerGroup.User

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
		userRouter.InitCommentRouter(MainGroup)
		userRouter.InitAdvertiseRouter(MainGroup)
		userRouter.InitMediaRouter(MainGroup)
		userRouter.InitMessagesRouter(MainGroup)
		userRouter.InitReportRouter(MainGroup)
		userRouter.InitNotificationRouter(MainGroup)
	}
	{
		adminRouter.InitAdminRouter(MainGroup)
		adminRouter.InitAdvertiseAdminRouter(MainGroup)
		adminRouter.InitRevenueAdminRouter(MainGroup)
		adminRouter.InitAdminReportRouter(MainGroup)
	}
	return r
}
