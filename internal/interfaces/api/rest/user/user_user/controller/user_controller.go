package controller

import "github.com/gin-gonic/gin"

type (
	IUserInfoController interface {
		GetInfoByUserId(ctx *gin.Context)
		GetManyUsers(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
	}
	IUserFriendController interface {
		SendAddFriendRequest(ctx *gin.Context)
		UndoFriendRequest(ctx *gin.Context)
		GetFriendRequests(ctx *gin.Context)
		AcceptFriendRequest(ctx *gin.Context)
		RejectFriendRequest(ctx *gin.Context)
		UnFriend(ctx *gin.Context)
		GetFriends(ctx *gin.Context)
		GetFriendSuggestion(ctx *gin.Context)
		GetFriendByBirthday(ctx *gin.Context)
		GetNonFriend(ctx *gin.Context)
	}
)
