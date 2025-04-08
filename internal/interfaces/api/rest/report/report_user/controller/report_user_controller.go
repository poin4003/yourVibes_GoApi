package controller

import "github.com/gin-gonic/gin"

type (
	IUserReportController interface {
		Report(ctx *gin.Context)
	}
)
