package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/controller"
	conversationRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type messagesRouter struct {
	conversationController       controller.IConversationController
	conversationDetailController controller.IConversationDetailController
	messageController            controller.IMessageController
	userProtectedMiddleware      middlewares.IUserAuthProtectedMiddleware
}

func NewMessagesRouter(
	conversationController controller.IConversationController,
	conversationDetailController controller.IConversationDetailController,
	messageController controller.IMessageController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *messagesRouter {
	return &messagesRouter{
		conversationController:       conversationController,
		conversationDetailController: conversationDetailController,
		messageController:            messageController,
		userProtectedMiddleware:      userProtectedMiddleware,
	}
}

func (r *messagesRouter) InitMessagesRouter(Router *gin.RouterGroup) {
	conversationRouter := Router.Group("/conversations")
	conversationRouter.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		conversationRouter.POST("/",
			helpers.ValidateFormBody(&conversationRequest.CreateConversationRequest{}, conversationRequest.ValidateCreateConversationRequest),
			r.conversationController.CreateConversation)

		conversationRouter.GET("/",
			helpers.ValidateQuery(&conversationQuery.ConversationObject{}, conversationQuery.ValidateConversationObject),
			r.conversationController.GetConversation,
		)

		conversationRouter.GET("/:conversationId", r.conversationController.GetConversationById)

		conversationRouter.DELETE("/:conversationId", r.conversationController.DeleteConversationById)

		conversationRouter.PATCH("/:conversationId",
			helpers.ValidateFormBody(&conversationRequest.UpdateConversationRequest{}, conversationRequest.ValidateUpdateConversationRequest),
			r.conversationController.UpdateConversation,
		)
	}

	messageRouterPublic := Router.Group("/messages")
	{
		messageRouterPublic.GET("/ws/:user_id", r.messageController.SendMessageWebSocket)
	}

	messageRouter := Router.Group("/messages")
	messageRouter.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		messageRouter.POST("/",
			helpers.ValidateJsonBody(&conversationRequest.CreateMessageRequest{}, conversationRequest.ValidateCreateMessageRequest),
			r.messageController.CreateMessage)

		messageRouter.GET("/get_by_conversation_id",
			helpers.ValidateQuery(&conversationQuery.MessageObject{}, conversationQuery.ValidateMessageQueryObject),
			r.messageController.GetMessagesByConversationId)

		messageRouter.GET("/message/:messageId", r.messageController.GetMessageById)

		messageRouter.DELETE("/message/:messageId", r.messageController.DeleteMessageById)
	}

	conversationDetailRouter := Router.Group("/conversation_details")
	conversationDetailRouter.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		conversationDetailRouter.POST("/",
			helpers.ValidateJsonBody(&conversationRequest.CreateConversationDetailRequest{}, conversationRequest.ValidateCreateConversationDetailRequest),
			r.conversationDetailController.CreateConversationDetail,
		)

		conversationDetailRouter.POST("/create_many",
			helpers.ValidateFormBody(&conversationRequest.CreateManyConversationDetailRequest{}, conversationRequest.ValidateCreateManyConversationDetailRequest),
			r.conversationDetailController.CreateManyConversationDetail,
		)

		conversationDetailRouter.GET("/get_by_id",
			helpers.ValidateQuery(&conversationQuery.ConversationDetailObject{}, conversationQuery.ValidateConversationDetailObject),
			r.conversationDetailController.GetConversationDetailByConversationId,
		)

		conversationDetailRouter.GET("/get_by_id/:userId/:conversationId", r.conversationDetailController.GetConversationDetailById)

		conversationDetailRouter.DELETE("/delete/:userId/:conversationId", r.conversationDetailController.DeleteConversationDetailById)

		conversationDetailRouter.PATCH("/update",
			helpers.ValidateJsonBody(&conversationRequest.UpdateConversationDetail{}, conversationRequest.ValidateUpdateConversationDetail),
			r.conversationDetailController.UpdateConversationDetail,
		)
	}
}
