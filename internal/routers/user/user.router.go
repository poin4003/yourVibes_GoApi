package user

import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// Public router

	//userRouterPublic := Router.Group("/user")
	//{
	//	userRouterPublic.POST("/register", account.Login.Register)
	//	userRouterPublic.POST("/login", account.Login.Login)
	//}

	// Private router
	//userRouterPrivate := Router.Group("/user")
	//userRouterPrivate.USE(limmiter())
	//userRouterPrivate.USE(Authen())
	//userRouterPrivate.Use(Permissions())
	//{
	//	userRouterPrivate.GET("/get_info")
	//}
}
