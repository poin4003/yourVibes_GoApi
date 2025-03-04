package message_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/message/service"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/message/message_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/message/message_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type MessageUserController struct {
	service message_service.MessageService
}

func NewMessageUserController(service message_service.MessageService) *MessageUserController {
	return &MessageUserController{service: service}
}

func (c *MessageUserController) CreateMessage(ctx *gin.Context) {
	var req request.CreateMessageRequest
	if err := ctx.BindJSON(&req); err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	msg, err := c.service.CreateMessage(ctx, req)
	if err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	extensions.ResponseSuccess(ctx, response.NewResponse(msg))
}

func (c *MessageUserController) GetMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		extensions.ResponseError(ctx, errors.New("id is required"))
		return
	}

	msg, err := c.service.GetMessage(ctx, id)
	if err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	extensions.ResponseSuccess(ctx, response.NewResponse(msg))
}

func (c *MessageUserController) UpdateMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		extensions.ResponseError(ctx, errors.New("id is required"))
		return
	}

	var req request.UpdateMessageRequest
	if err := ctx.BindJSON(&req); err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	msg, err := c.service.UpdateMessage(ctx, id, req)
	if err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	extensions.ResponseSuccess(ctx, response.NewResponse(msg))
}

func (c *MessageUserController) DeleteMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		extensions.ResponseError(ctx, errors.New("id is required"))
		return
	}

	err := c.service.DeleteMessage(ctx, id)
	if err != nil {
		extensions.ResponseError(ctx, err)
		return
	}

	extensions.ResponseSuccess(ctx, response.NewResponse(nil))
}