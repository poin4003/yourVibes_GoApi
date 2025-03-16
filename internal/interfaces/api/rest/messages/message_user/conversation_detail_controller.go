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

type cConversationController struct{}

func NewConversationDetailController() *cConversationController {
	return &cConversationController{}
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
func (cc *cConversationController) CreateConversationDetail(ctx *gin.Context) {
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
	createConversationDetailCommand, err := createConversationDetail.ToCreateConversationDetailCommand(createConversationDetail.UserId, createConversationDetail.ConversationId)

	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.ConversationDetail().CreateConversationDetail(ctx, createConversationDetailCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDetailDto := response.ToConversationDetailDto(result.ConversationDetail)

	pkgResponse.OK(ctx, conversationDetailDto)
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
func (cc *cConversationController) GetConversationDetailById(ctx *gin.Context) {
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

	result, err := services.ConversationDetail().GetConversationDetailById(ctx, userID, conversationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDetailDto := response.ToConversationDetailDto(result)

	pkgResponse.OK(ctx, conversationDetailDto)
}

// GetConversationDetailByUserId documentation
// @Summary Get conversationDetail by User ID
// @Description Retrieve a conversationDetail by its unique User ID
// @Tags conversationDetail
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param conversation_id query string false "Conversation ID"
// @Param limit query int false "Limit on page"
// @Param page query int false "Page number"
// @Security ApiKeyAuth
// @Router /conversation_details/get_by_id [get]
func (cc *cConversationController) GetConversationDetailByUserId(ctx *gin.Context) {
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

	if conversationDetailQueryObject.UserId == "" && conversationDetailQueryObject.ConversationId == "" {

		userIdClaim, err := extensions.GetUserID(ctx)
		if err != nil {
			ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
			return
		}
		conversationDetailQueryObject.UserId = userIdClaim.String()
	}

	getConversationDetailByUserIdQuery, _ := conversationDetailQueryObject.ToGetConversationDetailQuery()

	result, err := services.ConversationDetail().GetConversationDetailByUsesId(ctx, getConversationDetailByUserIdQuery)
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
func (cc *cConversationController) DeleteConversationDetailById(ctx *gin.Context) {
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
	err = services.ConversationDetail().DeleteConversationDetailById(ctx, deleteConversationDetailCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}
