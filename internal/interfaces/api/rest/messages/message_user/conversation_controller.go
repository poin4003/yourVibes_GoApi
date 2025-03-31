package message_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type cConversation struct {
}

func NewConversationController() *cConversation {
	return &cConversation{}
}

// CreateConversation documentation
// @Summary Conversation create Conversation
// @Description When user create conversation
// @Tags conversation
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false "Name of the conversation"
// @Param image formData file false "Image of the conversation" multiple
// @Param user_ids formData []string true "List of user IDs" collectionFormat(multi)
// @Security ApiKeyAuth
// @Router /conversations/ [post]
func (c *cConversation) CreateConversation(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}

	creatConsersation, ok := body.(*request.CreateConversationRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	var userIds []uuid.UUID
	for _, userId := range creatConsersation.UserIds {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			ctx.Error(pkgResponse.NewValidateError("Invalid user id"))
		}
		userIds = append(userIds, userUUID)
	}

	createConversation := creatConsersation.ToCreateConversationCommand(creatConsersation.Name, userIds)

	result, err := services.Conversation().CreateConversation(ctx, createConversation)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDto := response.ToConversationDto(result.Conversation)

	pkgResponse.OK(ctx, conversationDto)

}

// GetConversationById documentation
// @Summary Get conversation by ID
// @Description Retrieve a conversation by its unique ID
// @Tags conversation
// @Accept json
// @Produce json
// @Param conversation_id path string true "Conversation ID"
// @Security ApiKeyAuth
// @Router /conversations/{conversation_id} [get]
func (c *cConversation) GetConversationById(ctx *gin.Context) {
	// 1. Get conversation id from param
	conversationIdStr := ctx.Param("conversationId")
	conversatonId, err := uuid.Parse(conversationIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	result, err := services.Conversation().GetConversationById(ctx, conversatonId)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDto := response.ToConversationDto(result)

	pkgResponse.OK(ctx, conversationDto)
}

// DeleteConversation documentation
// @Summary delete conversation by ID
// @Description when user want to delete conversation
// @Tags conversation
// @Accept json
// @Produce json
// @Param conversation_id path string true "conversation ID"
// @Security ApiKeyAuth
// @Router /conversations/{conversation_id} [delete]
func (c *cConversation) DeleteConversationById(ctx *gin.Context) {

	// 1. Get conversation id from param
	conversationIdStr := ctx.Param("conversationId")
	conversationId, err := uuid.Parse(conversationIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	deleteConversationCommand := &command.DeleteConversationCommand{ConversationId: &conversationId}
	err = services.Conversation().DeleteConversationById(ctx, deleteConversationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// GetConversation documentation
// @Summary Get many conversation
// @Description When user get many conversation
// @Tags conversation
// @Accept json
// @Produce json
// @Param name query string false "Name of the conversation"
// @Param created_at query string false "Created at"
// @Param sort_by query string false "Sort by"
// @Param isDescending query bool false "Is descending"
// @Param limit query int false "Limit of conversation per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /conversations/ [get]
func (c *cConversation) GetConversation(ctx *gin.Context) {
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validatedQuery request"))
		return
	}

	userIdClaims, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	conversationQueryObject, ok := queryInput.(*query.ConversationObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid query type"))
		return
	}

	getManyConversationQuery, _ := conversationQueryObject.ToGetManyConversationQuery()

	result, err := services.Conversation().GetManyConversation(ctx, userIdClaims, getManyConversationQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	var conversationDtos []*response.ConversationDto
	for _, conversationResults := range result.Conversation {
		conversationDtos = append(conversationDtos, response.ToConversationDto(conversationResults))
	}

	pkgResponse.OKWithPaging(ctx, conversationDtos, *result.PagingResponse)

}

// UpdateConversation documentation
// @Summary update conversation
// @Description When user need to update conversation
// @Tags conversation
// @Accept multipart/form-data
// @Produce json
// @Param conversation_id path string true "ConversationId"
// @Param name formData string false "Name of the conversation"
// @Param image formData file false "Image of the conversation" multiple
// @Security ApiKeyAuth
// @Router /conversations/{conversation_id} [patch]
func (c *cConversation) UpdateConversation(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}

	updateConversationRequest, ok := body.(*request.UpdateConversationRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid update request type"))
	}

	conversationIdStr := ctx.Param("conversationId")
	conversationId, err := uuid.Parse(conversationIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	updateConversationCommand, _ := updateConversationRequest.ToUpdateConversationCommand(conversationId, &updateConversationRequest.Image)
	result, err := services.Conversation().UpdateConversationById(ctx, updateConversationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}
	conversationDto := response.ToConversationDto(result.Conversation)
	pkgResponse.OK(ctx, conversationDto)

}
