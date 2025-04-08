package controller

import "github.com/gin-gonic/gin"

type (
	IUserNotificationController interface {
		SendNotification(ctx *gin.Context)
		GetNotification(ctx *gin.Context)
		UpdateOneStatusNotifications(ctx *gin.Context)
		UpdateManyStatusNotifications(ctx *gin.Context)
	}
)
