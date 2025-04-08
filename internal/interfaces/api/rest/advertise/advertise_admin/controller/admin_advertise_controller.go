package controller

import "github.com/gin-gonic/gin"

type (
	IAdminAdvertiseController interface {
		GetAdvertiseDetail(ctx *gin.Context)
		GetManyAdvertise(ctx *gin.Context)
	}
)
