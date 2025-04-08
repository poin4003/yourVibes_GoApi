package controller

import "github.com/gin-gonic/gin"

type (
	ICommentLikeController interface {
		LikeComment(ctx *gin.Context)
		GetUserLikeComment(ctx *gin.Context)
	}
	ICommentUserController interface {
		CreateComment(ctx *gin.Context)
		GetComment(ctx *gin.Context)
		DeleteComment(ctx *gin.Context)
		UpdateComment(ctx *gin.Context)
	}
)
