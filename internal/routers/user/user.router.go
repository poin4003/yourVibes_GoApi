package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/user_auth"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/user_info"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	UserAuthController := user_auth.NewUserAuthController()
	UserInfoController := user_info.NewUserInfoController()
	// Public router

	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", UserAuthController.Register)
		userRouterPublic.POST("/verifyemail", UserAuthController.VerifyEmail)
		userRouterPublic.POST("/login", UserAuthController.Login)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(authentication.AuthProteced())
	{
		userRouterPrivate.GET("/:userId", UserInfoController.GetInfoByUserId)
		userRouterPrivate.GET("/", UserInfoController.GetManyUsers)
		userRouterPrivate.PATCH("/", UserInfoController.UpdateUser)
	}
}
