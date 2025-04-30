package controller

import "github.com/gin-gonic/gin"

type (
	IConversationController interface {
		CreateConversation(ctx *gin.Context)
		GetConversationById(ctx *gin.Context)
		DeleteConversationById(ctx *gin.Context)
		GetConversation(ctx *gin.Context)
		UpdateConversation(ctx *gin.Context)
	}
	IConversationDetailController interface {
		CreateConversationDetail(ctx *gin.Context)
		GetConversationDetailByConversationId(ctx *gin.Context)
		GetConversationDetailById(ctx *gin.Context)
		DeleteConversationDetailById(ctx *gin.Context)
		UpdateConversationDetail(ctx *gin.Context)
		CreateManyConversationDetail(ctx *gin.Context)
		TransferOwnerRole(ctx *gin.Context)
	}
	IMessageController interface {
		SendMessageWebSocket(ctx *gin.Context)
		CreateMessage(ctx *gin.Context)
		GetMessageById(ctx *gin.Context)
		GetMessagesByConversationId(ctx *gin.Context)
		DeleteMessageById(ctx *gin.Context)
	}
)
