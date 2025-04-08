package controller

import "github.com/gin-gonic/gin"

type (
	IAdminReportController interface {
		GetReportDetail(ctx *gin.Context)
		GetManyReports(ctx *gin.Context)
		HandleReport(ctx *gin.Context)
		DeleteReport(ctx *gin.Context)
		Activate(ctx *gin.Context)
	}
)
