package post_user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/ [post]
func (p *cPostUser) CreatePost(ctx *gin.Context) {
	var postInput request.CreatePostInput

	if err := ctx.ShouldBind(&postInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	if postInput.Content == "" && postInput.Media == nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "You must provide at least one of Content or Media")
		return
	}

	files := postInput.Media

	// Convert multipart.FileHeader to multipart.File
	var uploadedFiles []multipart.File
	for _, file := range files {
		openFile, err := file.Open()
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		uploadedFiles = append(uploadedFiles, openFile)
	}

	fmt.Println("Files retrieved:", len(files))

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	postModel := mapper.MapToPostFromCreateDto(&postInput, userUUID)
	post, resultCode, httpStatusCode, err := services.PostUser().CreatePost(context.Background(), postModel, uploadedFiles)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postDto := mapper.MapPostToNewPostDto(post)

	cacheKey := fmt.Sprintf("posts:user:%s:*", userUUID)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postDto)
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
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [patch]
func (p *cPostUser) UpdatePost(ctx *gin.Context) {
	var updateInput request.UpdatePostInput

	if err := ctx.ShouldBind(&updateInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	postFound, resultCodePostFound, httpStatusCodePostFound, err := services.PostUser().GetPost(ctx, postId, userIdClaim)
	if err != nil {
		response.ErrorResponse(ctx, resultCodePostFound, httpStatusCodePostFound, err.Error())
		return
	}

	if userIdClaim != postFound.UserId {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, fmt.Sprintf("You can not edit this post"))
		return
	}

	updateData := mapper.MapToPostFromUpdateDto(&updateInput)

	deleteMediaIds := updateInput.MediaIDs

	var uploadedFiles []multipart.File
	for _, fileHeader := range updateInput.Media {
		openFile, err := fileHeader.Open()
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		uploadedFiles = append(uploadedFiles, openFile)
	}

	post, resultCode, httpStatusCode, err := services.PostUser().UpdatePost(ctx, postId, updateData, deleteMediaIds, uploadedFiles)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postDto := mapper.MapPostToUpdatedPostDto(post)

	// Delete cache
	cacheKey := fmt.Sprintf("posts:user:%s:*", postFound.UserId)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postDto)
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
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /posts/ [get]
func (p *cPostUser) GetManyPost(ctx *gin.Context) {
	var query query.PostQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := fmt.Sprintf("posts:user:%s:page:%d:limit:%d", query.UserID, query.Page, query.Limit)
	cachePosts, err := p.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var postDto []response2.PostDto
		err = json.Unmarshal([]byte(cachePosts), &postDto)
		if err == nil {
			cacheTotalKey := fmt.Sprintf("posts:user:%s:total", query.UserID)
			cacheTatal, _ := p.redisClient.Get(context.Background(), cacheTotalKey).Int64()

			paging := response.PagingResponse{
				Limit: query.Limit,
				Page:  query.Page,
				Total: cacheTatal,
			}

			response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto, paging)
			return
		}
	}

	userUUID, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	postDtos, resultCode, httpStatusCode, paging, err := services.PostUser().GetManyPosts(ctx, &query, userUUID)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postsJson, _ := json.Marshal(postDtos)
	p.redisClient.Set(context.Background(), cacheKey, postsJson, time.Minute*1)

	cacheTotalKey := fmt.Sprintf("posts:user:%s:total", query.UserID)
	p.redisClient.Set(context.Background(), cacheTotalKey, paging.Total, time.Minute*1)

	response.SuccessPagingResponse(ctx, resultCode, http.StatusOK, postDtos, *paging)
}

// GetPostById documentation
// @Summary Get post by ID
// @Description Retrieve a post by its unique ID
// @Tags post_user
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [get]
func (p *cPostUser) GetPostById(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")

	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	cachedPost, err := p.redisClient.Get(context.Background(), postId.String()).Result()
	if err == nil {
		var postDto response2.PostDto
		err = json.Unmarshal([]byte(cachedPost), &postDto)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto)
		return
	}

	postDto, resultCode, httpStatusCode, err := services.PostUser().GetPost(ctx, postId, userIdClaim)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postJson, _ := json.Marshal(postDto)
	p.redisClient.Set(context.Background(), postId.String(), postJson, time.Minute*1)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postDto)
}

// DeletePost documentation
// @Summary delete post by ID
// @Description when user want to delete post
// @Tags post_user
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [delete]
func (p *cPostUser) DeletePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	postFound, resultCodePostFound, httpStatusCodePostFound, err := services.PostUser().GetPost(ctx, postId, userIdClaim)
	if err != nil {
		response.ErrorResponse(ctx, resultCodePostFound, httpStatusCodePostFound, err.Error())
		return
	}

	if userIdClaim != postFound.UserId {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, fmt.Sprintf("You can not delete this post"))
		return
	}

	resultCode, httpStatusCode, err := services.PostUser().DeletePost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	// Delete cache in redis
	cacheKey := fmt.Sprintf("posts:user:%s:*", postFound.UserId)
	keys, _, err := p.redisClient.Scan(ctx, 0, cacheKey, 0).Result()

	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	for _, key := range keys {
		if er := p.redisClient.Del(context.Background(), key).Err(); er != nil {
			panic(er.Error())
		}
	}

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postId)
}
