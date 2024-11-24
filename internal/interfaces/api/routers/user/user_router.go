package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth"
	auth_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user"
	user_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	UserAuthController := user_auth.NewUserAuthController()
	UserInfoController := user_user.NewUserInfoController()
	UserNotificationController := user_user.NewNotificationController()
	UserFriendController := user_user.NewUserFriendController()

	// Public router

	userRouterPublic := Router.Group("/users")
	{
		// user_auth
		userRouterPublic.POST("/register",
			helpers.ValidateJsonBody(&auth_request.RegisterRequest{}, auth_request.ValidateRegisterRequest),
			UserAuthController.Register,
		)

		userRouterPublic.POST("/verifyemail",
			helpers.ValidateJsonBody(&auth_request.VerifyEmailRequest{}, auth_request.ValidateVerifyEmailRequest),
			UserAuthController.VerifyEmail,
		)

		userRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&auth_request.LoginRequest{}, auth_request.ValidateLoginRequest),
			UserAuthController.Login,
		)

		// user_notification
		userRouterPublic.GET("/notifications/ws/:user_id", UserNotificationController.SendNotification)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middlewares.AuthProteced())
	{
		// user_info
		userRouterPrivate.GET("/:userId", UserInfoController.GetInfoByUserId)

		userRouterPrivate.GET("/",
			helpers.ValidateQuery(&user_query.UserQueryObject{}, user_query.ValidateUserQueryObject),
			UserInfoController.GetManyUsers,
		)

		userRouterPrivate.PATCH("/",
			helpers.ValidateFormBody(&user_request.UpdateUserRequest{}, user_request.ValidateUpdateUserRequest),
			UserInfoController.UpdateUser,
		)

		// user_notification
		userRouterPrivate.GET("/notifications",
			helpers.ValidateQuery(&user_query.NotificationQueryObject{}, user_query.ValidateNotificationQueryObject),
			UserNotificationController.GetNotification,
		)

		userRouterPrivate.PATCH("/notifications/:notification_id", UserNotificationController.UpdateOneStatusNotifications)
		userRouterPrivate.PATCH("/notifications", UserNotificationController.UpdateManyStatusNotifications)

		// user_friend
		userRouterPrivate.POST("/friends/friend_request/:friend_id", UserFriendController.SendAddFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_request/:friend_id", UserFriendController.UndoFriendRequest)

		userRouterPrivate.GET("/friends/friend_request",
			helpers.ValidateQuery(&user_query.FriendRequestQueryObject{}, user_query.ValidateFriendRequestQueryObject),
			UserFriendController.GetFriendRequests,
		)

		userRouterPrivate.POST("/friends/friend_response/:friend_id", UserFriendController.AcceptFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_response/:friend_id", UserFriendController.RejectFriendRequest)
		userRouterPrivate.DELETE("/friends/:friend_id", UserFriendController.UnFriend)

		userRouterPrivate.GET("/friends/:user_id",
			helpers.ValidateQuery(&user_query.FriendQueryObject{}, user_query.ValidateFriendQueryObject),
			UserFriendController.GetFriends,
		)
	}
}
