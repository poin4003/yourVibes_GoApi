package controller

import "github.com/gin-gonic/gin"

type (
	IPostLikeController interface {
		LikePost(ctx *gin.Context)
		GetUserLikePost(ctx *gin.Context)
	}
	IPostNewFeedController interface {
		DeleteNewFeed(ctx *gin.Context)
		GetNewFeeds(ctx *gin.Context)
	}
	IPostShareController interface {
		SharePost(ctx *gin.Context)
	}
	IPostUserController interface {
		CreatePost(ctx *gin.Context)
		UpdatePost(ctx *gin.Context)
		GetManyPost(ctx *gin.Context)
		GetPostById(ctx *gin.Context)
		DeletePost(ctx *gin.Context)
	}
)
