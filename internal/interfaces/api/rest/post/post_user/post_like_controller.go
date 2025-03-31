package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
)

type PostLikeController struct {
}

func NewPostLikeController() *PostLikeController {
	return &PostLikeController{}
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
		ctx.Error(response2.NewValidateError(err.Error()))
		return
	}

	// 2. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service to handle like or dislike
	likePostCommand := &command.LikePostCommand{PostId: postId, UserId: userIdClaim}

	result, err := services.LikeUserPost().LikePost(ctx, likePostCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	postDto := response.ToPostWithLikedDto(*result.Post)

	response2.OK(ctx, postDto)
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
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(response2.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	postLikeQueryObject, ok := queryInput.(*query.PostLikeQueryObject)
	if !ok {
		ctx.Error(response2.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(response2.NewValidateError("Invalid postId format"))
		return
	}

	// 4. Call service to get list user
	getPostLikeQuery, err := postLikeQueryObject.ToGetPostLikeQuery(postId)
	if err != nil {
		ctx.Error(response2.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.LikeUserPost().GetUsersOnLikes(ctx, getPostLikeQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var userDtos []*response.UserDto
	for _, userResult := range result.Users {
		userDtos = append(userDtos, response.ToUserDto(userResult))
	}

	response2.OKWithPaging(ctx, userDtos, *result.PagingResponse)
}
