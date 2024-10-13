package user_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/auth_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserAuth struct {
}

var UserAuth = new(cUserAuth)

// User verify email documentation
// @Summary User verify email
// @Description Before user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth_dto.VerifyEmailInput true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/verifyemail/ [post]
func (c *cUserAuth) VerifyEmail(ctx *gin.Context) {
	var input auth_dto.VerifyEmailInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidateParamEmail, http.StatusBadRequest, err.Error())
		return
	}

	code, err := services.UserAuth().VerifyEmail(ctx, input.Email)
	if err != nil {
		response.ErrorResponse(ctx, code, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(ctx, code, http.StatusOK, nil)
}

// User  documentation
// @Summary User Registration
// @Description When user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth_dto.RegisterCredentials true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/register/ [post]
func (c *cUserAuth) Register(ctx *gin.Context) {
	var registerInput auth_dto.RegisterCredentials

	if err := ctx.ShouldBindJSON(&registerInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidateParamRegister, http.StatusBadRequest, err.Error())
		return
	}

	code, err := services.UserAuth().Register(ctx, &registerInput)
	if err != nil {
		response.ErrorResponse(ctx, code, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(ctx, code, http.StatusOK, nil)
}

// User login documentation
// @Summary User login
// @Description When user login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth_dto.LoginCredentials true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Router /users/login/ [post]
func (c *cUserAuth) Login(ctx *gin.Context) {
	var loginInput auth_dto.LoginCredentials

	if err := ctx.ShouldBindJSON(&loginInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidateParamLogin, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, userModel, err := services.UserAuth().Login(ctx, &loginInput)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeLoginFailed, http.StatusBadRequest, err.Error())
		return
	}

	userDTO := mapper.MapUserToUserDto(userModel)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, gin.H{
		"access_token": accessToken,
		"user":         userDTO,
	})
}
