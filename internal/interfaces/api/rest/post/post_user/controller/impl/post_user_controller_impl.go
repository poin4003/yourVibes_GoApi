package impl

import (
	"context"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
)

type cPostUser struct {
	postService services.IPostUser
}

func NewPostUserController(
	postService services.IPostUser,
) *cPostUser {
	return &cPostUser{
		postService: postService,
	}
}

// CreatePost documentation
// @Summary Post create post
// @Description When user create post
// @Tags post_user
// @Accept multipart/form-data
// @Produce json
// @Param content formData string false "Content of the post"
// @Param privacy formData string false "Privacy level"
// @Param location formData string false "Location of the post"
// @Param media formData file false "Media files for the post" multiple
// @Security ApiKeyAuth
// @Router /posts/ [post]
func (c *cPostUser) CreatePost(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to updateUserRequest
	createPostRequest, ok := body.(*request.CreatePostRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle create post
	createPostCommand, err := createPostRequest.ToCreatePostCommand(userIdClaim, createPostRequest.Media)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = c.postService.CreatePost(context.Background(), createPostCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// UpdatePost documentation
// @Summary update post
// @Description When user need to update information of post or update media
// @Tags post_user
// @Accept multipart/form-data
// @Produce json
// @Param post_id path string true "PostId"
// @Param content formData string false "Post content"
// @Param privacy formData string false "Post privacy"
// @Param location formData string false "Post location"
// @Param media_ids formData int false "Array of mediaIds you want to delete"
// @Param media formData file false "Array of media you want to upload"
// @Security ApiKeyAuth
// @Router /posts/{post_id} [patch]
func (c *cPostUser) UpdatePost(ctx *gin.Context) {
	var postRequest query.PostQueryObject
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to updateUserRequest
	updatePostRequest, ok := body.(*request.UpdatePostRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get post_id from params
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 4. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 5. Call service to check owner
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	// 6. Get post to check owner
	queryResult, err := c.postService.GetPost(ctx, getOnePostQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 7. Check post advertise privacy
	if queryResult.Post.IsAdvertisement == consts.IS_ADVERTISE {
		if updatePostRequest.Privacy != pointer.Ptr(consts.PUBLIC) {
			ctx.Error(pkgResponse.NewCustomError(pkgResponse.ErrAdMustBePublic, "You can't update privacy of advertise"))
			return
		}
	}

	// 7. Get user id from token
	if userIdClaim != queryResult.Post.UserId {
		ctx.Error(pkgResponse.NewInvalidTokenError("You can not edit this post"))
		return
	}

	// 8. Call service to handle update post
	updatePostCommand, err := updatePostRequest.ToUpdatePostCommand(&postId, updatePostRequest.Media)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.postService.UpdatePost(ctx, updatePostCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 9. Map to dto
	postDto := response.ToPostDto(*result.Post)

	pkgResponse.OK(ctx, postDto)
}

// GetManyPost documentation
// @Summary Get many posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags post_user
// @Accept json
// @Produce json
// @Param user_id query string false "User ID to filter posts"
// @Param content query string false "Filter by content"
// @Param location query string false "Filter by location"
// @Param is_advertisement query int false "Filter by advertisement"
// @Param created_at query string false "Filter by creation time"
// @Param sort_by query string false "Which column to sort by"
// @Param isDescending query boolean false "Order by descending if true"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /posts/ [get]
func (c *cPostUser) GetManyPost(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to PostQueryObject
	postQueryObject, ok := queryInput.(*query.PostQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle get many
	getManyPostQuery, _ := postQueryObject.ToGetManyPostQuery(userIdClaim)

	result, err := c.postService.GetManyPosts(ctx, getManyPostQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var postDtos []*response.PostWithLikedDto
	for _, postResult := range result.Posts {
		postDtos = append(postDtos, response.ToPostWithLikedDto(*postResult))
	}

	pkgResponse.OKWithPaging(ctx, postDtos, *result.PagingResponse)
}

// GetTrendingPost documentation
// @Summary Get trending posts
// @Description Retrieve multiple trending posts
// @Tags post_user
// @Accept json
// @Produce json
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /posts/trending [get]
func (c *cPostUser) GetTrendingPost(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to TrendingPostQueryObject
	trendingPostQueryObject, ok := queryInput.(*query.TrendingPostQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle get many
	getTrendingPostQuery, _ := trendingPostQueryObject.ToGetTrendingQuery(userIdClaim)

	result, err := c.postService.GetTrendingPost(ctx, getTrendingPostQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var postDtos []*response.PostWithLikedDto
	for _, postResult := range result.Posts {
		postDtos = append(postDtos, response.ToPostWithLikedDto(*postResult))
	}

	pkgResponse.OKWithPaging(ctx, postDtos, *result.PagingResponse)
}

// GetPostById documentation
// @Summary Get post by ID
// @Description Retrieve a post by its unique ID
// @Tags post_user
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Security ApiKeyAuth
// @Router /posts/{post_id} [get]
func (c *cPostUser) GetPostById(ctx *gin.Context) {
	var postRequest query.PostQueryObject

	// 1. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
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

	// 3. Call service to handle get one
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}
	result, err := c.postService.GetPost(ctx, getOnePostQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to Dto
	postDto := response.ToPostWithLikedDto(*result.Post)

	pkgResponse.OK(ctx, postDto)
}

// DeletePost documentation
// @Summary delete post by ID
// @Description when user want to delete post
// @Tags post_user
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Security ApiKeyAuth
// @Router /posts/{post_id} [delete]
func (c *cPostUser) DeletePost(ctx *gin.Context) {
	var postRequest query.PostQueryObject

	// 1. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
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

	// 3. Get post to check owner
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	query_result, err := c.postService.GetPost(ctx, getOnePostQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Check owner
	if userIdClaim != query_result.Post.UserId {
		ctx.Error(pkgResponse.NewCustomError(pkgResponse.ErrInvalidToken, "you can not delete this post"))
		return
	}

	// 4. Call service delete
	deletePostCommand := &command.DeletePostCommand{PostId: &postId}

	err = c.postService.DeletePost(ctx, deletePostCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, postId)
}
