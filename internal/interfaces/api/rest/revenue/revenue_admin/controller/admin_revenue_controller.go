package controller

import "github.com/gin-gonic/gin"

type (
	IRevenueAdminController interface {
		GetMonthlyRevenue(ctx *gin.Context)
		GetSystemStats(ctx *gin.Context)
	}
)
