package controller

import "github.com/gin-gonic/gin"

type (
	IAdvertiseController interface {
		CreateAdvertise(ctx *gin.Context)
		GetManyAdvertise(ctx *gin.Context)
		GetAdvertiseWithStatistic(ctx *gin.Context)
	}
	IBillController interface {
		ConfirmPayment(ctx *gin.Context)
	}
)
