package user_auth

import (
	"github.com/gin-gonic/gin"
	userCommand "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserAuth struct {
}

func NewUserAuthController() *cUserAuth {
	return &cUserAuth{}
}

// VerifyEmail documentation
// @Summary User verify email
// @Description Before user registration
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.VerifyEmailRequest true "input"
// @Router /users/verifyemail/ [post]
func (c *cUserAuth) VerifyEmail(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to verify email request
	verifyEmailRequest, ok := body.(*request.VerifyEmailRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle verify email
	code, err := services.UserAuth().VerifyEmail(ctx, verifyEmailRequest.Email)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, code, http.StatusBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, code, http.StatusOK, nil)
}

// Register documentation
// @Summary User Registration
// @Description When user registration
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.RegisterRequest true "input"
// @Router /users/register/ [post]
func (c *cUserAuth) Register(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	registerRequest, ok := body.(*request.RegisterRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle register
	registerCommand, err := registerRequest.ToRegisterCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().Register(ctx, registerCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, nil)
}

// Login documentation
// @Summary User login
// @Description When user login
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.LoginRequest true "input"
// @Router /users/login/ [post]
func (c *cUserAuth) Login(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to loginRequest
	loginRequest, ok := body.(*request.LoginRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid login request type")
		return
	}

	// 3. Call service to handle login
	loginCommand, err := loginRequest.ToLoginCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().Login(ctx, loginCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Convert to dto
	userDto := response.ToUserWithSettingDto(result.User)

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"user":         userDto,
	})
}

// AuthGoogle documentation
// @Summary User auth google
// @Description When user need google login
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.AuthGoogleRequest true "input"
// @Router /users/auth_google/ [post]
func (c *cUserAuth) AuthGoogle(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to authGoogleRequest
	authGoogleRequest, ok := body.(*request.AuthGoogleRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle auth google
	authGoogleCommand, err := authGoogleRequest.ToAuthGoogleCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().AuthGoogle(ctx, authGoogleCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Map to dto
	userDto := response.ToUserWithSettingDto(result.User)

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"user":         userDto,
	})
}

// ChangePassword documentation
// @Summary User change password
// @Description When user need to change password
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.ChangePasswordRequest true "input"
// @Security ApiKeyAuth
// @Router /users/change_password/ [patch]
func (c *cUserAuth) ChangePassword(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to change password request
	changePasswordRequest, ok := body.(*request.ChangePasswordRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to handle change password
	changePasswordCommand, err := changePasswordRequest.ToChangePasswordCommand(userIdClaim)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().ChangePassword(ctx, changePasswordCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, nil)
}

// GetOtpForgotUserPassword documentation
// @Summary User get otp forgot user password
// @Description Before forgot password
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.GetOtpForgotUserPasswordRequest true "input"
// @Router /users/get_otp_forgot_user_password/ [post]
func (c *cUserAuth) GetOtpForgotUserPassword(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to get otp forgot user password request
	getOtpForgotUserPasswordRequest, ok := body.(*request.GetOtpForgotUserPasswordRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to get otp forgot user password request
	getOtpForgotUserPasswordCommand := &userCommand.GetOtpForgotUserPasswordCommand{
		Email: getOtpForgotUserPasswordRequest.Email,
	}

	result, err := services.UserAuth().GetOtpForgotUserPassword(ctx, getOtpForgotUserPasswordCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}

// ForgotUserPassword documentation
// @Summary User forgot password
// @Description When user forgot and change password
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.ForgotUserPasswordRequest true "input"
// @Router /users/forgot_user_password/ [post]
func (c *cUserAuth) ForgotUserPassword(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to forgot user password request
	forgotUserPasswordRequest, ok := body.(*request.ForgotUserPasswordRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle forgot user password
	forgotUserPasswordCommand, err := forgotUserPasswordRequest.ToForgotUserPasswordCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().ForgotUserPassword(ctx, forgotUserPasswordCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}

// AppAuthGoogle documentation
// @Summary User app auth google
// @Description When user need google login on mobile app
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.AppAuthGoogleRequest true "input"
// @Router /users/app_auth_google/ [post]
func (c *cUserAuth) AppAuthGoogle(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to appAuthGoogleRequest
	authGoogleRequest, ok := body.(*request.AppAuthGoogleRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle auth google
	appAuthGoogleCommand, err := authGoogleRequest.ToAppAuthGoogleCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().AppAuthGoogle(ctx, appAuthGoogleCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Map to dto
	userDto := response.ToUserWithSettingDto(result.User)

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"user":         userDto,
	})
}
