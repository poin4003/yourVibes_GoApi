package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
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
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [post]
func (p *PostLikeController) LikePost(ctx *gin.Context) {
	// 1. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 3. Call service to handle like or dislike
	likePostCommand := &command.LikePostCommand{PostId: postId, UserId: userIdClaim}

	result, err := services.LikeUserPost().LikePost(ctx, likePostCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	postDto := response.ToPostWithLikedDto(*result.Post)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, postDto)
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
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [get]
func (p *PostLikeController) GetUserLikePost(ctx *gin.Context) {
	// 1. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid postId format")
		return
	}

	// 2. Get query
	var query query.PostLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// 3. Call service to get list user
	getPostLikeQuery, err := query.ToGetPostLikeQuery(postId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.LikeUserPost().GetUsersOnLikes(ctx, getPostLikeQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var userDtos []*response.UserDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserDto(userResult))
	}

	pkg_response.SuccessPagingResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, userDtos, *result.PagingResponse)
}
