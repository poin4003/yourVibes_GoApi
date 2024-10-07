package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/user_auth"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// Public router

	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", user_auth.UserAuth.Register)
		userRouterPublic.POST("/verifyemail", user_auth.UserAuth.VerifyEmail)
	}

	// Private router
	//userRouterPrivate := Router.Group("/user")
	//userRouterPrivate.USE(limmiter())
	//userRouterPrivate.USE(Authen())
	//userRouterPrivate.Use(Permissions())
	//{
	//	userRouterPrivate.GET("/get_info")
	//}
}
