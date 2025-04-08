package controller

import "github.com/gin-gonic/gin"

type (
	IAdminAuthController interface {
		Login(ctx *gin.Context)
		ChangeAdminPassword(ctx *gin.Context)
	}
)
