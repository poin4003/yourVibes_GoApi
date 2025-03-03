package comment_user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cCommentLike struct{}

func NewCommentLikeController() *cCommentLike {
	return &cCommentLike{}
}

// LikeComment documentation
// @Summary Like comment
// @Description When user like comment
// @Tags like_comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID to create like comment"
// @Security ApiKeyAuth
// @Router /comments/like_comment/{comment_id} [post]
func (p *cCommentLike) LikeComment(ctx *gin.Context) {
	// 1. Get comment id form param
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service to handle like or dislike
	likeCommentCommand := &command.LikeCommentCommand{CommentId: commentId, UserId: userIdClaim}
	result, err := services.CommentLike().LikeComment(ctx, likeCommentCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	commentDto := response.ToCommentWithLikedDto(result.Comment)

	pkgResponse.OK(ctx, commentDto)
}

// GetUserLikeComment documentation
// @Summary Get User like comments
// @Description Retrieve multiple user is like comment
// @Tags like_comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID to get user like comment"
// @Param limit query int false "Limit of users per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /comments/like_comment/{comment_id} [get]
func (p *cCommentLike) GetUserLikeComment(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to CommentLikeQueryObject
	commentLikeQueryObject, ok := queryInput.(*query.CommentLikeQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get comment id from param
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(fmt.Sprintf("invalid comment id: %s", commentIdStr)))
		return
	}

	// 4. Call service to handle get user like comment
	getUserLikeCommentQuery, err := commentLikeQueryObject.ToGetCommentLikeQuery(commentId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.CommentLike().GetUsersOnLikeComment(ctx, getUserLikeCommentQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var userDtos []*response.UserDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserDto(userResult))
	}

	pkgResponse.OKWithPaging(ctx, userDtos, *result.PagingResponse)
}
