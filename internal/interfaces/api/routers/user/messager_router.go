package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user"
	conversationRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type MessagesRouter struct{}

func (mr *MessagesRouter) InitMessagesRouter(Router *gin.RouterGroup) {
	userConversationController := message_user.NewConversationController()
	useMessageController := message_user.NewMessageController()
	useConversationDetailController := message_user.NewConversationDetailController()

	conversationRouter := Router.Group("/conversations")
	{
		conversationRouter.POST("/create_conversation",
			helpers.ValidateJsonBody(&conversationRequest.CreateConversationRequest{}, conversationRequest.ValidateCreateConversationRequest),
			userConversationController.CreateConversation)

		conversationRouter.GET("/conversation",
			helpers.ValidateQuery(&conversationQuery.ConversationObject{}, conversationQuery.ValidateConversationObject),
			userConversationController.GetConversation)

		conversationRouter.GET("/conversation/:conversationId", userConversationController.GetConversationById)

		conversationRouter.DELETE("/conversation/:conversationId", userConversationController.DeleteConversationById)

	}

	messageRouter := Router.Group("/messages")
	messageRouter.Use(middlewares.UserAuthProtected())
	{
		messageRouter.POST("/create_message",
			helpers.ValidateJsonBody(&conversationRequest.CreateMessageRequest{}, conversationRequest.ValidateCreateMessageRequest),
			useMessageController.CreateMessage)

		messageRouter.GET("/get_by_conversation_id",
			helpers.ValidateQuery(&conversationQuery.MessageObject{}, conversationQuery.ValidateMessageQueryObject),
			useMessageController.GetMessagesByConversationId)

		messageRouter.GET("/message/:messageId", useMessageController.GetMessageById)
	}

	conversationDetailRouter := Router.Group("/conversation_detail")
	conversationDetailRouter.Use(middlewares.UserAuthProtected())
	{
		conversationDetailRouter.POST("/create_conversation_detail",
			helpers.ValidateJsonBody(&conversationRequest.CreateConversationDetailRequest{}, conversationRequest.ValidateCreatCOnversationDetailRequest),
			useConversationDetailController.CreateConversationDetail)
		conversationDetailRouter.GET("/get_by_id/:userId/:conversationId", useConversationDetailController.GetConversationDetailById)
		conversationDetailRouter.GET("/get_by_user_id",
			helpers.ValidateQuery(&conversationQuery.ConversationDetailObject{}, conversationQuery.ValidateConversationDetailObject),
			useConversationDetailController.GetConversationDetailByUserId)
	}

}
