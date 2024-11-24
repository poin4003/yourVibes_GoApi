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
// @Router /users/register/ [post]
func (c *cUserAuth) Register(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	registerRequest, ok := body.(*request.RegisterRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle register
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
// @Router /users/login/ [post]
func (c *cUserAuth) Login(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to loginRequest
	loginRequest, ok := body.(*request.LoginRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid login request type")
		return
	}

	// 3. Call service to handle login
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

	// 4. Convert to dto
	userDto := response.ToUserWithSettingDto(result.User)

	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"user":         userDto,
	})
}
