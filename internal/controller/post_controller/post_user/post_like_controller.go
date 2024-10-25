package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/user_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Tags Like_Post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to create like post"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [post]
func (p *PostLikeController) LikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}
	likeUserPostModel := mapper.MapToLikeUserPostFromPostIdAndUserId(postId, userUUID)
	postModel, resultCode, httpStatusCode, err := services.LikeUserPost().LikePost(ctx, likeUserPostModel)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postDto := mapper.MapPostToPostDto(postModel)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, httpStatusCode, postDto)
}

// GetUserLikePost documentation
// @Summary Get User like posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags Like_Post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID to get user like post"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/like_post/{post_id} [get]
func (p *PostLikeController) GetUserLikePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "Invalid postId format")
		return
	}

	var query query_object.PostLikeQueryObject
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	likeUserPost, resultCode, httpStatusCode, paging, err := services.LikeUserPost().GetUsersOnLikes(ctx, postId, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	var userDtos []user_dto.UserDtoShortVer
	for _, user := range likeUserPost {
		userDto := mapper.MapUserToUserDtoShortVer(user)
		userDtos = append(userDtos, userDto)
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, userDtos, *paging)
}
