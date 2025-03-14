package message_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type cMessage struct {
}

func NewMessageController() *cMessage {
	return &cMessage{}
}

// CreateMessage documentation
// @Summary Message create Message
// @Description When user create message
// @Tags message
// @Accept json
// @Produce json
// @Param input body request.CreateMessageRequest true "input"
// @Security ApiKeyAuth
// @Router /messages/ [post]
func (m *cMessage) CreateMessage(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}

	createMessageRequest, ok := body.(*request.CreateMessageRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	userIdClaims, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	createMessageCommand, err := createMessageRequest.ToCreateMessageCommand(userIdClaims)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.Message().CreateMessage(ctx, createMessageCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	messageDto := response.ToMessageDto(result.Message)

	pkgResponse.OK(ctx, messageDto)

}

// GetMessageById documentation
// @Summary Get message by ID
// @Description Retrieve a message by its unique ID
// @Tags message
// @Accept json
// @Produce json
// @Param messageId path string true "Message ID"
// @Security ApiKeyAuth
// @Router /messages/message/{messageId} [get]
func (m *cMessage) GetMessageById(ctx *gin.Context) {
	messageIdStr := ctx.Param("messageId")
	messageId, err := uuid.Parse(messageIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	result, err := services.Message().GetMessageById(ctx, messageId)
	if err != nil {
		ctx.Error(err)
		return
	}

	messageDto := response.ToMessageDto(result)

	pkgResponse.OK(ctx, messageDto)
}

// GetMessagesByConversationId documentation
// @Summary Get messages by conversation ID
// @Description Retrieve messages by conversation ID
// @Tags message
// @Accept json
// @Produce json
// @Param conversation_id query string true "ConversationId"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Security ApiKeyAuth
// @Router /messages/get_by_conversation_id [get]
func (m *cMessage) GetMessagesByConversationId(ctx *gin.Context) {
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewValidateError("Missing validateQuery request"))
		return
	}

	MessagesByConversationIdQuery, ok := queryInput.(*query.MessageObject)
	if !ok {
		ctx.Error(pkgResponse.NewValidateError("Invalid query type"))
		return
	}

	getMessagesByConversationIdQuery, _ := MessagesByConversationIdQuery.ToGetManyMessageQuery()

	result, err := services.Message().GetMessagesByConversationId(ctx, getMessagesByConversationIdQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	var messageDtos []*response.MessageDto
	for _, messageResult := range result.Messages {
		messageDtos = append(messageDtos, response.ToMessageDto(messageResult))
	}

	pkgResponse.OKWithPaging(ctx, messageDtos, *result.PagingResponse)
}
