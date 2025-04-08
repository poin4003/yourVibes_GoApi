package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
)

type cCommentUser struct {
	commentUserService services.ICommentUser
}

func NewCommentUserController(
	commentUserService services.ICommentUser,
) *cCommentUser {
	return &cCommentUser{
		commentUserService: commentUserService,
	}
}

// CreateComment documentation
// @Summary Comment create comment
// @Description When user create comment or rep comment
// @Tags comment_user
// @Accept json
// @Produce json
// @Param input body request.CreateCommentRequest true "input"
// @Security ApiKeyAuth
// @Router /comments/ [post]
func (c *cCommentUser) CreateComment(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to create comment request
	createCommentRequest, ok := body.(*request.CreateCommentRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 2. Get userid from token
	userIdClaims, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service to handle create comment
	createCommentCommand, err := createCommentRequest.ToCreateCommentCommand(userIdClaims)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.commentUserService.CreateComment(ctx, createCommentCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	commentDto := response.ToCommentDto(result.Comment)

	pkgResponse.OK(ctx, commentDto)
}

// GetComment documentation
// @Summary Get many comment
// @Description Retrieve multiple comment filtered by various criteria.
// @Tags comment_user
// @Accept json
// @Produce json
// @Param post_id query string true "Post ID to filter comment, get the first layer"
// @Param parent_id query string false "Filter by parent id, get the next layer"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /comments/ [get]
func (c *cCommentUser) GetComment(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to CommentQueryObject
	commentQueryObject, ok := queryInput.(*query.CommentQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 2. Get user id from token
	userIdClaims, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service to handle get many
	getManyCommentQuery, err := commentQueryObject.ToGetManyCommentQuery(userIdClaims)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.commentUserService.GetManyComments(ctx, getManyCommentQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	var commentDtos []*response.CommentWithLikedDto
	for _, commentResult := range result.Comments {
		commentDtos = append(commentDtos, response.ToCommentWithLikedDto(commentResult))
	}

	pkgResponse.OKWithPaging(ctx, commentDtos, *result.PagingResponse)
}

// DeleteComment documentation
// @Summary delete comment by ID
// @Description when user want to delete comment
// @Tags comment_user
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID"
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [delete]
func (c *cCommentUser) DeleteComment(ctx *gin.Context) {
	// 1. Get comment id form param
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Call service to handle delete comment
	deleteCommentCommand := &command.DeleteCommentCommand{CommentId: commentId}
	err = c.commentUserService.DeleteComment(ctx, deleteCommentCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// UpdateComment documentation
// @Summary update comment
// @Description When user need to update information of comment
// @Tags comment_user
// @Accept json
// @Produce json
// @Param comment_id path string true "commentId"
// @Param input body request.UpdateCommentRequest true "input"
// @Security ApiKeyAuth
// @Router /comments/{comment_id} [patch]
func (c *cCommentUser) UpdateComment(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to update comment request
	updateCommentRequest, ok := body.(*request.UpdateCommentRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 2. Get commend id from param
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service to handle update comment
	updateCommentCommand, err := updateCommentRequest.ToUpdateCommentCommand(commentId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.commentUserService.UpdateComment(ctx, updateCommentCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	commentDto := response.ToCommentDto(result.Comment)

	pkgResponse.OK(ctx, commentDto)
}
