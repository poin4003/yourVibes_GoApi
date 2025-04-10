package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	cUserAuth "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/controller"
	authRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/request"
	cUser "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/controller"
	userRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
)

type userRouter struct {
	userInfoController      cUser.IUserInfoController
	userFriendController    cUser.IUserFriendController
	userAuthController      cUserAuth.IUserAuthController
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware
}

func NewUserRouter(
	userController cUser.IUserInfoController,
	userFriendController cUser.IUserFriendController,
	userAuthController cUserAuth.IUserAuthController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *userRouter {
	return &userRouter{
		userInfoController:      userController,
		userFriendController:    userFriendController,
		userAuthController:      userAuthController,
		userProtectedMiddleware: userProtectedMiddleware,
	}
}

func (r *userRouter) InitUserRouter(Router *gin.RouterGroup) {
	// Public router
	userRouterPublic := Router.Group("/users")
	{
		// user_auth
		userRouterPublic.POST("/register",
			helpers.ValidateJsonBody(&authRequest.RegisterRequest{}, authRequest.ValidateRegisterRequest),
			r.userAuthController.Register,
		)

		userRouterPublic.POST("/verifyemail",
			helpers.ValidateJsonBody(&authRequest.VerifyEmailRequest{}, authRequest.ValidateVerifyEmailRequest),
			r.userAuthController.VerifyEmail,
		)

		userRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&authRequest.LoginRequest{}, authRequest.ValidateLoginRequest),
			r.userAuthController.Login,
		)

		userRouterPublic.POST("/app_auth_google",
			helpers.ValidateJsonBody(&authRequest.AppAuthGoogleRequest{}, authRequest.ValidateAppAuthGoogleRequest),
			r.userAuthController.AppAuthGoogle,
		)

		userRouterPublic.POST("/auth_google",
			helpers.ValidateJsonBody(&authRequest.AuthGoogleRequest{}, authRequest.ValidateAuthGoogleRequest),
			r.userAuthController.AuthGoogle,
		)

		userRouterPublic.POST("/get_otp_forgot_user_password",
			helpers.ValidateJsonBody(&authRequest.GetOtpForgotUserPasswordRequest{}, authRequest.ValidateGetOtpForgotUserPasswordRequest),
			r.userAuthController.GetOtpForgotUserPassword,
		)

		userRouterPublic.POST("/forgot_user_password",
			helpers.ValidateJsonBody(&authRequest.ForgotUserPasswordRequest{}, authRequest.ValidateForgotUserPasswordRequest),
			r.userAuthController.ForgotUserPassword,
		)
	}

	// Private router
	userRouterPrivate := Router.Group("/users")
	userRouterPrivate.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		// user authentication
		userRouterPrivate.PATCH("/change_password",
			helpers.ValidateJsonBody(&authRequest.ChangePasswordRequest{}, authRequest.ValidateChangePasswordRequest),
			r.userAuthController.ChangePassword,
		)

		// user_info
		userRouterPrivate.GET("/:userId", r.userInfoController.GetInfoByUserId)

		userRouterPrivate.GET("/",
			helpers.ValidateQuery(&userQuery.UserQueryObject{}, userQuery.ValidateUserQueryObject),
			r.userInfoController.GetManyUsers,
		)

		userRouterPrivate.PATCH("/",
			helpers.ValidateFormBody(&userRequest.UpdateUserRequest{}, userRequest.ValidateUpdateUserRequest),
			r.userInfoController.UpdateUser,
		)

		// user_friend
		userRouterPrivate.POST("/friends/friend_request/:friend_id", r.userFriendController.SendAddFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_request/:friend_id", r.userFriendController.UndoFriendRequest)

		userRouterPrivate.GET("/friends/friend_request",
			helpers.ValidateQuery(&userQuery.FriendRequestQueryObject{}, userQuery.ValidateFriendRequestQueryObject),
			r.userFriendController.GetFriendRequests,
		)

		userRouterPrivate.POST("/friends/friend_response/:friend_id", r.userFriendController.AcceptFriendRequest)
		userRouterPrivate.DELETE("/friends/friend_response/:friend_id", r.userFriendController.RejectFriendRequest)
		userRouterPrivate.DELETE("/friends/:friend_id", r.userFriendController.UnFriend)

		userRouterPrivate.GET("/friends/:user_id",
			helpers.ValidateQuery(&userQuery.FriendQueryObject{}, userQuery.ValidateFriendQueryObject),
			r.userFriendController.GetFriends,
		)

		userRouterPrivate.GET("/friends/suggestion",
			helpers.ValidateQuery(&userQuery.FriendQueryObject{}, userQuery.ValidateFriendQueryObject),
			r.userFriendController.GetFriendSuggestion,
		)

		userRouterPrivate.GET("/friends/birthday",
			helpers.ValidateQuery(&userQuery.FriendQueryObject{}, userQuery.ValidateFriendQueryObject),
			r.userFriendController.GetFriendByBirthday,
		)

		userRouterPrivate.GET("/friends/non_friend",
			helpers.ValidateQuery(&userQuery.FriendQueryObject{}, userQuery.ValidateFriendQueryObject),
			r.userFriendController.GetNonFriend,
		)
	}
}
