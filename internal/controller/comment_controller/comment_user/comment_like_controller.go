package comment_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cCommentLike struct{}

func NewCommentLikeController() *cCommentLike {
	return &cCommentLike{}
}

// LikeComment documentation
// @Summary Like comment
// @Description When user like comment
// @Tags Like_Comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID to create like comment"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /comments/like_comment/{comment_id} [post]
func (p *cCommentLike) LikeComment(ctx *gin.Context) {
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	likeUserCommentModel := mapper.MapToLikeUserCommentFromCommentIdAndUserId(commentId, userUUID)

	resultCode, httpStatusCode, err := services.CommentLike().LikeComment(ctx, likeUserCommentModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, httpStatusCode, nil)
}

// GetUserLikeComment documentation
// @Summary Get User like comments
// @Description Retrieve multiple user is like comment
// @Tags Like_Comment
// @Accept json
// @Produce json
// @Param comment_id path string true "comment ID to get user like comment"
// @Param limit query int false "Limit of users per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/like_comment/{comment_id} [get]
func (p *cCommentLike) GetUserLikeComment(ctx *gin.Context) {
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, fmt.Sprintf("invalid comment id: %s", commentIdStr))
		return
	}

	var query query_object.CommentLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, fmt.Sprintf("invalid query"))
		return
	}

	likeUserComment, resultCode, httpStatusCode, paging, err := services.CommentLike().GetUsersOnLikeComment(ctx, commentId, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	var userDtos []user_dto.UserDtoShortVer
	for _, user := range likeUserComment {
		userDto := mapper.MapUserToUserDtoShortVer(user)
		userDtos = append(userDtos, userDto)
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, userDtos, *paging)
}
