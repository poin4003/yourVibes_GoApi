package comment_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
	user_mapper "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
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
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	likeUserCommentModel := mapper.MapToLikeUserCommentFromCommentIdAndUserId(commentId, userUUID)

	commentDto, resultCode, httpStatusCode, err := services.CommentLike().LikeComment(ctx, likeUserCommentModel, userUUID)
	if err != nil {
		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, commentDto)
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
	commentIdStr := ctx.Param("comment_id")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, fmt.Sprintf("invalid comment id: %s", commentIdStr))
		return
	}

	var query query.CommentLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, fmt.Sprintf("invalid query"))
		return
	}

	likeUserComment, resultCode, httpStatusCode, paging, err := services.CommentLike().GetUsersOnLikeComment(ctx, commentId, &query)
	if err != nil {
		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	var userDtos []response.UserDtoShortVer
	for _, user := range likeUserComment {
		userDto := user_mapper.MapUserToUserDtoShortVer(user)
		userDtos = append(userDtos, userDto)
	}

	pkg_response.SuccessPagingResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, userDtos, *paging)
}
