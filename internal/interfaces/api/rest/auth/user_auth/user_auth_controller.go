package user_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/verifyemail/ [post]
func (c *cUserAuth) VerifyEmail(ctx *gin.Context) {
	var input request.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidateParamEmail, http.StatusBadRequest, err.Error())
		return
	}

	code, err := services.UserAuth().VerifyEmail(ctx, input.Email)
	if err != nil {
		pkg_response.ErrorResponse(ctx, code, http.StatusBadRequest, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, code, http.StatusOK, nil)
}

// Register documentation
// @Summary User Registration
// @Description When user registration
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.RegisterRequest true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/register/ [post]
func (c *cUserAuth) Register(ctx *gin.Context) {
	var registerRequest request.RegisterRequest

	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidateParamRegister, http.StatusBadRequest, err.Error())
		return
	}

	registerCommand, err := registerRequest.ToRegisterCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().Register(ctx, registerCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, nil)
}

// Login documentation
// @Summary User login
// @Description When user login
// @Tags user_auth
// @Accept json
// @Produce json
// @Param input body request.LoginRequest true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/login/ [post]
func (c *cUserAuth) Login(ctx *gin.Context) {
	var loginRequest request.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidateParamLogin, http.StatusBadRequest, err.Error())
		return
	}

	loginCommand, err := loginRequest.ToLoginCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserAuth().Login(ctx, loginCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeLoginFailed, http.StatusBadRequest, err.Error())
		return
	}

	userDto := response.ToUserWithSettingDto(result.User)

	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"user":         userDto,
	})
}
