package message_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cConversation struct {
}

func NewConversationController() *cConversation {
	return &cConversation{}
}

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

	createConversation := creatConsersation.ToCreateConversationCommand(creatConsersation.Name)

	result, err := services.Conversation().CreateConversation(ctx, createConversation)
	if err != nil {
		ctx.Error(err)
		return
	}

	conversationDto := response.ToConversationDto(result.Conversation)

	pkgResponse.OK(ctx, conversationDto)

}

func (c *cConversation) GetConversation(ctx *gin.Context) {
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validatedQuery request"))
		return
	}

	conversationQueryObject, ok := queryInput.(*query.ConversationObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid query type"))
		return
	}

	getManyConversationQuery, _ := conversationQueryObject.ToGetManyConversationQuery()

	result, err := services.Conversation().GetManyConversation(ctx, getManyConversationQuery)
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

func (c *cConversation) GetConversationById(ctx *gin.Context) {
	// var conversatonRequest query.ConversationObject

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
