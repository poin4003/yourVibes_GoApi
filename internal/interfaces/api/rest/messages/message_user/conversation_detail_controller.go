package message_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cConversationController struct{}

func NewConversationDetailController() *cConversationController {
	return &cConversationController{}
}

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

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	getConversationDetailByUserIdQuery, _ := conversationDetailQueryObject.ToGetConversationDetailQuery(userIdClaim, conversationDetailQueryObject.ConversationId)

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
