package comment_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/comment_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cCommentUser struct {
}

func NewCommentUserController() *cCommentUser {
	return &cCommentUser{}
}

func (p *cCommentUser) CreateComment(ctx *gin.Context) {
	var commentInput comment_dto.CreateCommentInput

	if err := ctx.ShouldBindJSON(&commentInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	commentModel := mapper.MapToCommentFromCreateDto(&commentInput, userUUID)
	comment, resultCode, err := services.CommentUser().CreateComment(ctx, commentModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, comment)
}

func (p *cCommentUser) GetComment(ctx *gin.Context) {
	var query query_object.CommentQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	comment, resultCode, err := services.CommentUser().GetManyComments(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	total := int64(len(comment))
	paging := response.PagingResponse{
		Limit: query.Limit,
		Page:  query.Page,
		Total: total,
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, comment, paging)
}
