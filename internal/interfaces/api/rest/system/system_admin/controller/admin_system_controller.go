package controller

import "github.com/gin-gonic/gin"

type (
	ISystemAdminCache interface {
		ClearAllCache(ctx *gin.Context)
		ClearAllPostCache(ctx *gin.Context)
		ClearAllCommentCache(ctx *gin.Context)
	}
	ISystemAdminPost interface {
		UpdatePostAndStatistics(ctx *gin.Context)
		DelayPostCreatedAt(ctx *gin.Context)
		ExpiredAdvertiseByPostId(ctx *gin.Context)
		PushAdvertiseToNewFeed(ctx *gin.Context)
		PushFeaturePostToNewFeed(ctx *gin.Context)
		CheckExpiryOfAdvertisement(ctx *gin.Context)
		CheckExpiryOfFeaturePost(ctx *gin.Context)
	}
)
