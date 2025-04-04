package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user"
	conversationRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type MessagesRouter struct{}

func (mr *MessagesRouter) InitMessagesRouter(Router *gin.RouterGroup) {
	userConversationController := message_user.NewConversationController()
	useMessageController := message_user.NewMessageController(global.MessageSocketHub)
	useConversationDetailController := message_user.NewConversationDetailController()

	conversationRouter := Router.Group("/conversations")
	conversationRouter.Use(middlewares.UserAuthProtected())
	{
		conversationRouter.POST("/",
			helpers.ValidateFormBody(&conversationRequest.CreateConversationRequest{}, conversationRequest.ValidateCreateConversationRequest),
			userConversationController.CreateConversation)

		conversationRouter.GET("/",
			helpers.ValidateQuery(&conversationQuery.ConversationObject{}, conversationQuery.ValidateConversationObject),
			userConversationController.GetConversation)

		conversationRouter.GET("/:conversationId", userConversationController.GetConversationById)

		conversationRouter.DELETE("/:conversationId", userConversationController.DeleteConversationById)

		conversationRouter.PATCH("/:conversationId",
			helpers.ValidateFormBody(&conversationRequest.UpdateConversationRequest{}, conversationRequest.ValidateUpdateConversationRequest),
			userConversationController.UpdateConversation)
	}

	messageRouterPublic := Router.Group("/messages")
	{
		messageRouterPublic.GET("/ws/:user_id", useMessageController.SendMessageWebSocket)
	}

	messageRouter := Router.Group("/messages")
	messageRouter.Use(middlewares.UserAuthProtected())
	{
		messageRouter.POST("/",
			helpers.ValidateJsonBody(&conversationRequest.CreateMessageRequest{}, conversationRequest.ValidateCreateMessageRequest),
			useMessageController.CreateMessage)

		messageRouter.GET("/get_by_conversation_id",
			helpers.ValidateQuery(&conversationQuery.MessageObject{}, conversationQuery.ValidateMessageQueryObject),
			useMessageController.GetMessagesByConversationId)

		messageRouter.GET("/message/:messageId", useMessageController.GetMessageById)

		messageRouter.DELETE("/message/:messageId", useMessageController.DeleteMessageById)
	}

	conversationDetailRouter := Router.Group("/conversation_details")
	conversationDetailRouter.Use(middlewares.UserAuthProtected())
	{
		conversationDetailRouter.POST("/",
			helpers.ValidateJsonBody(&conversationRequest.CreateConversationDetailRequest{}, conversationRequest.ValidateCreatCOnversationDetailRequest),
			useConversationDetailController.CreateConversationDetail)

		conversationDetailRouter.GET("/get_by_id",
			helpers.ValidateQuery(&conversationQuery.ConversationDetailObject{}, conversationQuery.ValidateConversationDetailObject),
			useConversationDetailController.GetConversationDetailByConversationId)

		conversationDetailRouter.GET("/get_by_id/:userId/:conversationId", useConversationDetailController.GetConversationDetailById)

		conversationDetailRouter.DELETE("/delete/:userId/:conversationId", useConversationDetailController.DeleteConversationDetailById)

		conversationDetailRouter.PATCH("/update",
			helpers.ValidateJsonBody(&conversationRequest.UpdateConversationDetail{}, conversationRequest.ValidateUpdateConversationDetail),
			useConversationDetailController.UpdateConversationDetail)
	}

}
