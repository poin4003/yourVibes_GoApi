package controller

import "github.com/gin-gonic/gin"

type (
	IUserAuthController interface {
		VerifyEmail(ctx *gin.Context)
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		AuthGoogle(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		GetOtpForgotUserPassword(ctx *gin.Context)
		ForgotUserPassword(ctx *gin.Context)
		AppAuthGoogle(ctx *gin.Context)
	}
)
