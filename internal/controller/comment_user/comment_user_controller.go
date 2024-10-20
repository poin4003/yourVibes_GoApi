package comment_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// CreateComment documentation
// @Summary Comment create comment
// @Description When user create comment or rep comment
// @Tags comment
// @Accept json
// @Produce json
// @Param input body comment_dto.CreateCommentInput true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/ [post]
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

	commentDto := mapper.MapCommentToCommentDto(comment)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, commentDto)
}

// GetManyComment documentation
// @Summary Get many comment
// @Description Retrieve multiple comment filtered by various criteria.
// @Tags comment
// @Accept json
// @Produce json
// @Param post_id query string true "Post ID to filter comment, get the first layer"
// @Param parent_id query string false "Filter by parent id, get the next layer"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/ [get]
func (p *cCommentUser) GetComment(ctx *gin.Context) {
	var query query_object.CommentQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	comments, resultCode, paging, err := services.CommentUser().GetManyComments(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	var commentDtos []comment_dto.CommentDto
	for _, comment := range comments {
		commentDto := mapper.MapCommentToCommentDto(comment)
		commentDtos = append(commentDtos, *commentDto)
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, commentDtos, *paging)
}

// DeleteComment documentation
// @Summary delete comment by ID
// @Description when user want to delete comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [delete]
func (p *cCommentUser) DeleteComment(ctx *gin.Context) {
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	resultCode, err := services.CommentUser().DeleteComment(ctx, commentId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, http.StatusOK)
}

// UpdateComment documentation
// @Summary update comment
// @Description When user need to update information of comment
// @Tags comment
// @Accept multipart/form-data
// @Produce json
// @Param comment_id path string true "commentId"
// @Param input body comment_dto.UpdateCommentInput true "input"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [patch]
func (p *cCommentUser) UpdateComment(ctx *gin.Context) {
	var updateInput comment_dto.UpdateCommentInput

	if err := ctx.ShouldBindJSON(&updateInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	commentModel := mapper.MapToCommentFromUpdateDto(&updateInput)

	comment, resultCode, err := services.CommentUser().UpdateComment(ctx, commentId, commentModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	commentDto := mapper.MapCommentToCommentDto(comment)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, commentDto)
}
