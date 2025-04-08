package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
)

type cConversationDetail struct {
	conversationDetailService services.IConversationDetail
}

func NewConversationDetailController(
	conversationDetailService services.IConversationDetail,
) *cConversationDetail {
	return &cConversationDetail{
		conversationDetailService: conversationDetailService,
	}
}

// CreateConversationDetail documentation
// @Summary ConversationDetail create ConversationDatail
// @Description When user create conversationDetail
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param input body request.CreateConversationDetailRequest true "input"
// @Security ApiKeyAuth
// @Router /conversation_details/ [post]
func (c *cConversationDetail) CreateConversationDetail(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}
	createConversationDetail, ok := body.(*request.CreateConversationDetailRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}
	createConversationDetailCommand, err := createConversationDetail.ToCreateConversationDetailCommand()

	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.conversationDetailService.CreateConversationDetail(ctx, createConversationDetailCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDetailDto := response.ToConversationDetailDto(result.ConversationDetail)

	pkgResponse.OK(ctx, conversationDetailDto)
}

// GetConversationDetailByConversationId documentation
// @Summary Get conversationDetail by Conversation Id response List User in Conversation
// @Description Retrieve a conversationDetail by its unique Conversation ID
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param conversation_id query string true "Conversation ID"
// @Param limit query int false "Limit on page"
// @Param page query int false "Page number"
// @Security ApiKeyAuth
// @Router /conversation_details/get_by_id [get]
func (c *cConversationDetail) GetConversationDetailByConversationId(ctx *gin.Context) {
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validatedQuery request"))
		return
	}

	conversationDetailQueryObject, ok := queryInput.(*query.ConversationDetailObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid query type"))
		return
	}

	getConversationDetailByUserIdQuery, _ := conversationDetailQueryObject.ToGetConversationDetailQuery()

	result, err := c.conversationDetailService.GetConversationDetailByConversationId(ctx, getConversationDetailByUserIdQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	var conversationDetailDtos []*response.ConversationDetailDto
	for _, conversationDetailResults := range result.ConversationDetail {
		conversationDetailDtos = append(conversationDetailDtos, response.ToConversationDetailDto(conversationDetailResults))
	}

	pkgResponse.OKWithPaging(ctx, conversationDetailDtos, *result.PagingResponse)
}

// GetConversationDetailById documentation
// @Summary Get conversationDetail by ID
// @Description Retrieve a conversationDetail by its unique ID
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param conversationId path string true "Conversation ID"
// @Security ApiKeyAuth
// @Router /conversation_details/get_by_id/{userId}/{conversationId} [get]
func (c *cConversationDetail) GetConversationDetailById(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	conversationIdStr := ctx.Param("conversationId")

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	conversationID, err := uuid.Parse(conversationIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.conversationDetailService.GetConversationDetailById(ctx, userID, conversationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDetailDto := response.ToConversationDetailDto(result)

	pkgResponse.OK(ctx, conversationDetailDto)
}

// DeleteConversationDetailById documentation
// @Summary Delete conversationDetail by ID
// @Description when user delete conversationDetail
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param conversation_id path string true "Conversation ID"
// @Security ApiKeyAuth
// @Router /conversation_details/delete/{user_id}/{conversation_id} [delete]
func (c *cConversationDetail) DeleteConversationDetailById(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	conversationIdStr := ctx.Param("conversationId")

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	conversationID, err := uuid.Parse(conversationIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	deleteConversationDetailCommand := &command.DeleteConversationDetailCommand{UserId: &userID, ConversationId: &conversationID}
	err = c.conversationDetailService.DeleteConversationDetailById(ctx, deleteConversationDetailCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// UpdateConversationDetail Update Last Message Status of notification to false
// @Summary Update conversationDetail status to false
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param input body request.UpdateConversationDetail true "input"
// @Security ApiKeyAuth
// @Router /conversation_details/update [patch]
func (c *cConversationDetail) UpdateConversationDetail(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}

	updateConversationDetail, ok := body.(*request.UpdateConversationDetail)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	updateOneStatusConversationDetailCommand, err := updateConversationDetail.ToUpdateConversationDetailCommand(
		updateConversationDetail.UserId,
		updateConversationDetail.ConversationId,
	)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = c.conversationDetailService.UpdateOneStatusConversationDetail(ctx, updateOneStatusConversationDetailCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// CreateManyConversationDetail documentation
// @Summary Create many conversationDetail
// @Description When user create many conversationDetail
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param input body request.CreateManyConversationDetailRequest true "input"
// @Security ApiKeyAuth
// @Router /conversation_details/create_many [post]
func (c *cConversationDetail) CreateManyConversationDetail(ctx *gin.Context) {
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validateRequest request"))
		return
	}

	createManyConversation, ok := body.(*request.CreateManyConversationDetailRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	var userIds []uuid.UUID
	for _, userId := range createManyConversation.UserIds {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			ctx.Error(pkgResponse.NewValidateError("Invalid user id"))
			return
		}
		userIds = append(userIds, userUUID)
	}

	createManyConversationCommand, err := createManyConversation.ToCreateManyConversationDetailCommands(
		userIds,
	)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}
	result, err := c.conversationDetailService.CreateManyConversationDetail(ctx, createManyConversationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}
	conversationDetailDtos := response.ToManyConversationDetailDto(result.ConversationDetails)

	pkgResponse.OK(ctx, conversationDetailDtos)
}
