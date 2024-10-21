package post_user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/dtos/post_dto"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/redis/go-redis/v9"
	"mime/multipart"
	"net/http"
	"time"
)

type PostUserController struct {
	redisClient *redis.Client
}

func NewPostUserController(
	redisClient *redis.Client,
) *PostUserController {
	return &PostUserController{
		redisClient: redisClient,
	}
}

var (
	validate = validator.New()
)

// CreatePost documentation
// @Summary Post create post
// @Description When user create post
// @Tags post
// @Accept multipart/form-data
// @Produce json
// @Param title formData string false "Title of the post"
// @Param content formData string false "Content of the post"
// @Param privacy formData string false "Privacy level"
// @Param location formData string false "Location of the post"
// @Param media formData file false "Media files for the post" multiple
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/ [post]
func (p *PostUserController) CreatePost(ctx *gin.Context) {
	var postInput post_dto.CreatePostInput

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
	post, resultCode, err := services.PostUser().CreatePost(context.Background(), postModel, uploadedFiles)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	postDto := mapper.MapPostToNewPostDto(post)

	cacheKey := fmt.Sprintf("posts:user:%s:page:*:limit:*", userUUID)
	p.redisClient.Del(context.Background(), cacheKey)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto)
}

// UpdatePost documentation
// @Summary update post
// @Description When user need to update information of post or update media
// @Tags post
// @Accept multipart/form-data
// @Produce json
// @Param post_id path string true "PostId"
// @Param title formData string false "Post title"
// @Param content formData string false "Post content"
// @Param privacy formData string false "Post privacy"
// @Param location formData string false "Post location"
// @Param media_ids formData int false "Array of mediaIds you want to delete"
// @Param media formData file false "Array of media you want to upload"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [patch]
func (p *PostUserController) UpdatePost(ctx *gin.Context) {
	var updateInput post_dto.UpdatePostInput

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

	postFound, resultCodePostFound, err := services.PostUser().GetPost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCodePostFound, http.StatusInternalServerError, err.Error())
		return
	}

	if postFound == nil {
		response.ErrorResponse(ctx, response.ErrDataNotFound, http.StatusBadRequest, fmt.Sprintf("post id %s not found", postIdStr))
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

	post, resultCode, err := services.PostUser().UpdatePost(ctx, postId, updateData, deleteMediaIds, uploadedFiles)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	postDto := mapper.MapPostToPostDto(post)

	// Delete cache
	cacheKey := fmt.Sprintf("posts:user:%s:page:*:limit:*", postFound.UserId)
	cachePostKey := postId.String()
	p.redisClient.Del(context.Background(), cacheKey)
	p.redisClient.Del(context.Background(), cachePostKey)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto)
}

// GetManyPost documentation
// @Summary Get many posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags post
// @Accept json
// @Produce json
// @Param user_id query string false "User ID to filter posts"
// @Param title query string false "Filter by post title"
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
func (p *PostUserController) GetManyPost(ctx *gin.Context) {
	var query query_object.PostQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	cacheKey := fmt.Sprintf("posts:user:%s:page:%d:limit:%d", query.UserID, query.Page, query.Limit)
	cachePosts, err := p.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var postDto []post_dto.PostDto
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

	posts, resultCode, paging, err := services.PostUser().GetManyPosts(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	var postDtos []post_dto.PostDto
	for _, post := range posts {
		postDto := mapper.MapPostToPostDto(post)
		postDtos = append(postDtos, *postDto)
	}

	postsJson, _ := json.Marshal(postDtos)
	p.redisClient.Set(context.Background(), cacheKey, postsJson, time.Minute*1)

	cacheTotalKey := fmt.Sprintf("posts:user:%s:total", query.UserID)
	p.redisClient.Set(context.Background(), cacheTotalKey, paging.Total, time.Minute*1)

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDtos, *paging)
}

// GetPostById documentation
// @Summary Get post by ID
// @Description Retrieve a post by its unique ID
// @Tags post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [get]
func (p *PostUserController) GetPostById(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")

	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	cachedPost, err := p.redisClient.Get(context.Background(), postId.String()).Result()
	if err == nil {
		var postDto post_dto.PostDto
		err = json.Unmarshal([]byte(cachedPost), &postDto)
		if err != nil {
			response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto)
		return
	}

	var post *model.Post
	var resultCode int

	post, resultCode, err = services.PostUser().GetPost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	postDto := mapper.MapPostToPostDto(post)

	postJson, _ := json.Marshal(postDto)
	p.redisClient.Set(context.Background(), postId.String(), postJson, time.Minute*1)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, postDto)
}

// DeletePost documentation
// @Summary delete post by ID
// @Description when user want to delete post
// @Tags post
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{post_id} [delete]
func (p *PostUserController) DeletePost(ctx *gin.Context) {
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

	postFound, resultCodePostFound, err := services.PostUser().GetPost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCodePostFound, http.StatusInternalServerError, err.Error())
		return
	}

	if postFound == nil {
		response.ErrorResponse(ctx, response.ErrDataNotFound, http.StatusBadRequest, fmt.Sprintf("post id %s not found", postIdStr))
		return
	}

	if userIdClaim != postFound.UserId {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusForbidden, fmt.Sprintf("You can not delete this post"))
		return
	}

	resultCode, err := services.PostUser().DeletePost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete cache in redis
	cacheKey := fmt.Sprintf("posts:user:%s:page:*:limit:*", postFound.UserId)
	cachePostKey := postId.String()
	p.redisClient.Del(context.Background(), cacheKey)
	p.redisClient.Del(context.Background(), cachePostKey)

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusNoContent, postId)
}
