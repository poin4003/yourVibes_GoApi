package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/controller/auth/user_auth"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/controller/user/user_user"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	UserAuthController := user_auth.NewUserAuthController()
	UserInfoController := user_user.NewUserInfoController()
	UserNotificationController := user_user.NewNotificationController()
	UserFriendController := user_user.NewUserFriendController()
	UserNewFeedController := user_user.NewUserNewFeedController()
	// Public router

	userRouterPublic := Router.Group("/users")
	{
		// user_auth
		userRouterPublic.POST("/register", UserAuthController.Register)
		userRouterPublic.POST("/verifyemail", UserAuthController.VerifyEmail)
		userRouterPublic.POST("/login", UserAuthController.Login)

		// user_notification
		userRouterPublic.GET("/notifications/ws/:user_id", UserNotificationController.SendNotification)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middlewares.AuthProteced())
	{
		// user_info
		userRouterPrivate.GET("/:userId", UserInfoController.GetInfoByUserId)
		userRouterPrivate.GET("/", UserInfoController.GetManyUsers)
		userRouterPrivate.PATCH("/", UserInfoController.UpdateUser)

		// user_notification
		userRouterPrivate.GET("/notifications", UserNotificationController.GetNotification)
		userRouterPrivate.PATCH("/notifications/:notification_id", UserNotificationController.UpdateOneStatusNotifications)
		userRouterPrivate.PATCH("/notifications", UserNotificationController.UpdateManyStatusNotifications)

		// user_friend
		userRouterPrivate.POST("/friends/friend_request/:friend_id", UserFriendController.SendAddFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_request/:friend_id", UserFriendController.UndoFriendRequest)
		userRouterPrivate.GET("/friends/friend_request", UserFriendController.GetFriendRequests)
		userRouterPrivate.POST("/friends/friend_response/:friend_id", UserFriendController.AcceptFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_response/:friend_id", UserFriendController.RejectFriendRequest)
		userRouterPrivate.DELETE("/friends/:friend_id", UserFriendController.UnFriend)
		userRouterPrivate.GET("/friends/", UserFriendController.GetFriends)

		// user_new_feed
		userRouterPrivate.DELETE("/new_feeds/:post_id", UserNewFeedController.DeleteNewFeed)
		userRouterPrivate.GET("/new_feeds/", UserNewFeedController.GetNewFeeds)
	}
}
