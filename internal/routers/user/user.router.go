package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/user_auth"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// Public router

	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", user_auth.UserAuth.Register)
		userRouterPublic.POST("/verifyemail", user_auth.UserAuth.VerifyEmail)
		userRouterPublic.POST("/login", user_auth.UserAuth.Login)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middlewares.AuthProteced())
	{
		userRouterPrivate.GET("/get_info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "Ok",
			})
		})
	}
}
