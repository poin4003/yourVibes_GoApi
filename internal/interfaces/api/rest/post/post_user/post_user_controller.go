package post_user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"github.com/redis/go-redis/v9"
	"mime/multipart"
	"net/http"
	"time"
)

type cPostUser struct {
	redisClient *redis.Client
}

func NewPostUserController(
	redisClient *redis.Client,
) *cPostUser {
	return &cPostUser{
		redisClient: redisClient,
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
func (p *cPostUser) CreatePost(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to updateUserRequest
	createPostRequest, ok := body.(*request.CreatePostRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get file from body
	files := createPostRequest.Media

	// 3.1. Convert multipart.FileHeader to multipart.File
	var uploadedFiles []multipart.File
	for _, file := range files {
		openFile, err := file.Open()
		if err != nil {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		uploadedFiles = append(uploadedFiles, openFile)
	}

	// 4. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 5. Call service to handle create post
	createPostCommand, err := createPostRequest.ToCreatePostCommand(userIdClaim, uploadedFiles)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostUser().CreatePost(context.Background(), createPostCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 6. Map result to dto
	postDto := response.ToPostDto(*result.Post)

	cacheKey := fmt.Sprintf("posts:user:%s:*", userIdClaim)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, postDto)
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
func (p *cPostUser) UpdatePost(ctx *gin.Context) {
	var postRequest query.PostQueryObject
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to updateUserRequest
	updatePostRequest, ok := body.(*request.UpdatePostRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get post_id from params
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 5. Call service to check owner
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Get post to check owner
	queryResult, err := services.PostUser().GetPost(ctx, getOnePostQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, queryResult.ResultCode, queryResult.HttpStatusCode, err.Error())
		return
	}

	// 7. Check post advertise privacy
	if queryResult.Post.IsAdvertisement {
		if updatePostRequest.Privacy != pointer.Ptr(consts.PUBLIC) {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrAdMustBePublic, http.StatusBadRequest, "You can't update privacy of advertise")
			return
		}
	}

	// 7. Get user id from token
	if userIdClaim != queryResult.Post.UserId {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusForbidden, fmt.Sprintf("You can not edit this post"))
		return
	}

	// 8. Get file from body
	var uploadedFiles []multipart.File
	for _, fileHeader := range updatePostRequest.Media {
		openFile, err := fileHeader.Open()
		if err != nil {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		uploadedFiles = append(uploadedFiles, openFile)
	}

	// 9. Call service to handle update post
	updatePostCommand, err := updatePostRequest.ToUpdatePostCommand(&postId, uploadedFiles)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostUser().UpdatePost(ctx, updatePostCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 10. Map to dto
	postDto := response.ToPostDto(*result.Post)

	// Delete cache
	cacheKey := fmt.Sprintf("posts:user:%s:*", result.Post.UserId)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, postDto)
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
// @Param is_advertisement query boolean false "Filter by advertisement"
// @Param created_at query string false "Filter by creation time"
// @Param sort_by query string false "Which column to sort by"
// @Param isDescending query boolean false "Order by descending if true"
// @Param limit query int false "Limit of posts per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /posts/ [get]
func (p *cPostUser) GetManyPost(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to PostQueryObject
	postQueryObject, ok := queryInput.(*query.PostQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 2. Check redis
	cacheKey := fmt.Sprintf("posts:user:%s:page:%d:limit:%d", postQueryObject.UserID, postQueryObject.Page, postQueryObject.Limit)
	cachePosts, err := p.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var postDto []response.PostDto
		err = json.Unmarshal([]byte(cachePosts), &postDto)
		if err == nil {
			cacheTotalKey := fmt.Sprintf("posts:user:%s:total", postQueryObject.UserID)
			cacheTatal, _ := p.redisClient.Get(context.Background(), cacheTotalKey).Int64()

			paging := pkg_response.PagingResponse{
				Limit: postQueryObject.Limit,
				Page:  postQueryObject.Page,
				Total: cacheTatal,
			}

			pkg_response.SuccessPagingResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, postDto, paging)
			return
		}
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to handle get many
	getManyPostQuery, err := postQueryObject.ToGetManyPostQuery(userIdClaim)

	result, err := services.PostUser().GetManyPosts(ctx, getManyPostQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map to dto
	var postDtos []*response.PostWithLikedDto
	for _, postResult := range result.Posts {
		postDtos = append(postDtos, response.ToPostWithLikedDto(*postResult))
	}

	postsJson, _ := json.Marshal(postDtos)
	p.redisClient.Set(context.Background(), cacheKey, postsJson, time.Minute*1)

	cacheTotalKey := fmt.Sprintf("posts:user:%s:total", postQueryObject.UserID)
	p.redisClient.Set(context.Background(), cacheTotalKey, result.PagingResponse.Total, time.Minute*1)

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, postDtos, *result.PagingResponse)
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
func (p *cPostUser) GetPostById(ctx *gin.Context) {
	var postRequest query.PostQueryObject

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

	// 3. Check redis
	cachedPost, err := p.redisClient.Get(context.Background(), postId.String()).Result()
	if err == nil {
		var postDto response.PostDto
		err = json.Unmarshal([]byte(cachedPost), &postDto)
		if err != nil {
			pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		pkg_response.SuccessResponse(ctx, pkg_response.ErrCodeSuccess, http.StatusOK, postDto)
		return
	}

	// 4. Call service to handle get one
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.PostUser().GetPost(ctx, getOnePostQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map to Dto
	postDto := response.ToPostWithLikedDto(*result.Post)

	postJson, _ := json.Marshal(postDto)
	p.redisClient.Set(context.Background(), postId.String(), postJson, time.Minute*1)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, postDto)
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
func (p *cPostUser) DeletePost(ctx *gin.Context) {
	var postRequest query.PostQueryObject

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

	// 3. Get post to check owner
	getOnePostQuery, err := postRequest.ToGetOnePostQuery(postId, userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	query_result, err := services.PostUser().GetPost(ctx, getOnePostQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, query_result.ResultCode, query_result.HttpStatusCode, err.Error())
		return
	}

	// 4. Check owner
	if userIdClaim != query_result.Post.UserId {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusForbidden, fmt.Sprintf("You can not delete this post"))
		return
	}

	// 4. Call service delete
	deletePostCommand := &command.DeletePostCommand{PostId: &postId}

	result, err := services.PostUser().DeletePost(ctx, deletePostCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// Delete cache in redis
	cacheKey := fmt.Sprintf("posts:user:%s:*", query_result.Post.UserId)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, postId)
}
