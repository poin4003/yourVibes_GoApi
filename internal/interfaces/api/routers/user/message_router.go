package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/message/message_user"
)

func MessageRouter(router *gin.RouterGroup) {
	messageRouter := router.Group("/messages")
	{
		messageRouter.GET("", message_user.GetMessages)
		messageRouter.POST("", message_user.CreateMessage)
		messageRouter.GET("/:id", message_user.GetMessage)
		messageRouter.PUT("/:id", message_user.UpdateMessage)
		messageRouter.DELETE("/:id", message_user.DeleteMessage)
	}
}
