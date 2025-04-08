package controller

import "github.com/gin-gonic/gin"

type (
	IAdminInfoController interface {
		UpdateAdminInfo(ctx *gin.Context)
	}
)
