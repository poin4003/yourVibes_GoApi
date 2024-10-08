package user_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cUserAuth struct {
}

var UserAuth = new(cUserAuth)

func (c *cUserAuth) VerifyEmail(ctx *gin.Context) {
	var input vo.VerifyEmailInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(ctx, 2, err.Error())
		return
	}

	code, err := services.UserAuth().VerifyEmail(ctx, input.Email)
	if err != nil {
		response.ErrorResponse(ctx, code, err.Error())
		return
	}

	response.SuccessResponse(ctx, code, nil)
}

func (c *cUserAuth) Register(ctx *gin.Context) {
	var registerInput vo.RegisterCredentials

	if err := ctx.ShouldBindJSON(&registerInput); err != nil {
		response.ErrorResponse(ctx, 2, err.Error())
		return
	}

	code, err := services.UserAuth().Register(ctx, &registerInput)
	if err != nil {
		response.ErrorResponse(ctx, code, err.Error())
		return
	}

	response.SuccessResponse(ctx, code, nil)
}

func (c *cUserAuth) Login(ctx *gin.Context) {
	var loginInput vo.LoginCredentials

	if err := ctx.ShouldBindJSON(&loginInput); err != nil {
		response.ErrorResponse(ctx, 2, err.Error())
		return
	}

	accessToken, user, err := services.UserAuth().Login(ctx, &loginInput)
	if err != nil {
		response.ErrorResponse(ctx, 3, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, gin.H{
		"access_token": accessToken,
		"user":         user,
	})
}
