package controller

import "github.com/gin-gonic/gin"

type (
	IMediaController interface {
		GetMedia(ctx *gin.Context)
	}
)
