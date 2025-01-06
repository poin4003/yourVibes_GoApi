package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth"
	authRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user"
	userRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	UserAuthController := user_auth.NewUserAuthController()
	UserInfoController := user_user.NewUserInfoController()
	UserNotificationController := user_user.NewNotificationController()
	UserFriendController := user_user.NewUserFriendController()
	UserReportController := user_user.NewUserReportController()

	// Public router

	userRouterPublic := Router.Group("/users")
	{
		// user_auth
		userRouterPublic.POST("/register",
			helpers.ValidateJsonBody(&authRequest.RegisterRequest{}, authRequest.ValidateRegisterRequest),
			UserAuthController.Register,
		)

		userRouterPublic.POST("/verifyemail",
			helpers.ValidateJsonBody(&authRequest.VerifyEmailRequest{}, authRequest.ValidateVerifyEmailRequest),
			UserAuthController.VerifyEmail,
		)

		userRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&authRequest.LoginRequest{}, authRequest.ValidateLoginRequest),
			UserAuthController.Login,
		)

		userRouterPublic.POST("/app_auth_google",
			helpers.ValidateJsonBody(&authRequest.AppAuthGoogleRequest{}, authRequest.ValidateAppAuthGoogleRequest),
			UserAuthController.AppAuthGoogle,
		)

		userRouterPublic.POST("/auth_google",
			helpers.ValidateJsonBody(&authRequest.AuthGoogleRequest{}, authRequest.ValidateAuthGoogleRequest),
			UserAuthController.AuthGoogle,
		)

		userRouterPublic.POST("/get_otp_forgot_user_password",
			helpers.ValidateJsonBody(&authRequest.GetOtpForgotUserPasswordRequest{}, authRequest.ValidateGetOtpForgotUserPasswordRequest),
			UserAuthController.GetOtpForgotUserPassword,
		)

		userRouterPublic.POST("/forgot_user_password",
			helpers.ValidateJsonBody(&authRequest.ForgotUserPasswordRequest{}, authRequest.ValidateForgotUserPasswordRequest),
			UserAuthController.ForgotUserPassword,
		)

		// user_notification
		userRouterPublic.GET("/notifications/ws/:user_id", UserNotificationController.SendNotification)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		// user authentication
		userRouterPrivate.PATCH("/change_password",
			helpers.ValidateJsonBody(&authRequest.ChangePasswordRequest{}, authRequest.ValidateChangePasswordRequest),
			UserAuthController.ChangePassword,
		)

		// user_info
		userRouterPrivate.GET("/:userId", UserInfoController.GetInfoByUserId)

		userRouterPrivate.GET("/",
			helpers.ValidateQuery(&userQuery.UserQueryObject{}, userQuery.ValidateUserQueryObject),
			UserInfoController.GetManyUsers,
		)

		userRouterPrivate.PATCH("/",
			helpers.ValidateFormBody(&userRequest.UpdateUserRequest{}, userRequest.ValidateUpdateUserRequest),
			UserInfoController.UpdateUser,
		)

		// user_notification
		userRouterPrivate.GET("/notifications",
			helpers.ValidateQuery(&userQuery.NotificationQueryObject{}, userQuery.ValidateNotificationQueryObject),
			UserNotificationController.GetNotification,
		)

		userRouterPrivate.PATCH("/notifications/:notification_id", UserNotificationController.UpdateOneStatusNotifications)
		userRouterPrivate.PATCH("/notifications", UserNotificationController.UpdateManyStatusNotifications)

		// user_friend
		userRouterPrivate.POST("/friends/friend_request/:friend_id", UserFriendController.SendAddFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_request/:friend_id", UserFriendController.UndoFriendRequest)

		userRouterPrivate.GET("/friends/friend_request",
			helpers.ValidateQuery(&userQuery.FriendRequestQueryObject{}, userQuery.ValidateFriendRequestQueryObject),
			UserFriendController.GetFriendRequests,
		)

		userRouterPrivate.POST("/friends/friend_response/:friend_id", UserFriendController.AcceptFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_response/:friend_id", UserFriendController.RejectFriendRequest)
		userRouterPrivate.DELETE("/friends/:friend_id", UserFriendController.UnFriend)

		userRouterPrivate.GET("/friends/:user_id",
			helpers.ValidateQuery(&userQuery.FriendQueryObject{}, userQuery.ValidateFriendQueryObject),
			UserFriendController.GetFriends,
		)

		// user report
		userRouterPrivate.POST("/report",
			helpers.ValidateJsonBody(&userRequest.ReportUserRequest{}, userRequest.ValidateReportUserRequest),
			UserReportController.ReportUser,
		)
	}
}
