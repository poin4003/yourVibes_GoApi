package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	mapper2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type PostLikeController struct {
	redisClient *redis.Client
}

func NewPostLikeController(
	redisClient *redis.Client,
) *PostLikeController {
	return &PostLikeController{
		redisClient: redisClient,
	}
}

// LikePost documentation
// @Summary Like Post
// @Description When user like post
// @Tags like_post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to create like post"
// @Success 200 {object} pkg_response.pkg_responseData
// @Failure 500 {object} pkg_response.Errpkg_response
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [post]
func (p *PostLikeController) LikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}
	likeUserPostModel := mapper2.MapToLikeUserPostFromPostIdAndUserId(postId, userIdClaim)
	postDto, resultCode, httpStatusCode, err := services.LikeUserPost().LikePost(ctx, likeUserPostModel, userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, httpStatusCode, postDto)
}

// GetUserLikePost documentation
// @Summary Get User like posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags like_post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to get user like post"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} pkg_response.pkg_responseData
// @Failure 500 {object} pkg_response.Errpkg_response "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [get]
func (p *PostLikeController) GetUserLikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid postId format")
		return
	}

	var query query.PostLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	likeUserPost, resultCode, httpStatusCode, paging, err := services.LikeUserPost().GetUsersOnLikes(ctx, postId, &query)
	if err != nil {
		pkg_response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	var userDtos []response.UserDtoShortVer
	for _, user := range likeUserPost {
		userDto := mapper.MapUserToUserDtoShortVer(user)
		userDtos = append(userDtos, userDto)
	}

	pkg_response.SuccessPagingResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, userDtos, *paging)
}
