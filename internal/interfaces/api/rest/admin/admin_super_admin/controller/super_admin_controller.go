package controller

import "github.com/gin-gonic/gin"

type (
	ISuperAdminController interface {
		CreateAdmin(ctx *gin.Context)
		UpdateAdmin(ctx *gin.Context)
		GetAdminById(ctx *gin.Context)
		GetManyAdmins(ctx *gin.Context)
		ForgotAdminPassword(ctx *gin.Context)
	}
)
