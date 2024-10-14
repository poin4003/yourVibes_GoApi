package post_user

import (
	"context"
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
	"mime/multipart"
	"net/http"
)

type PostUserController struct{}

func NewPostUserController() *PostUserController {
	return &PostUserController{}
}

var PostUser = new(PostUserController)

var (
	validate = validator.New()
)

// User create post documentation
// @Summary Post create post
// @Description When user create post
// @Tags post
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title of the post"
// @Param content formData string true "Content of the post"
// @Param privacy formData string true "Privacy level"
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

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, post)
}

// Update post documentation
// @Summary update post
// @Description When user need to update information of post or update media
// @Tags post
// @Accept multipart/form-data
// @Produce json
// @Param postId path string true "PostId"
// @Param title formData string false "Post title"
// @Param content formData string false "Post content"
// @Param privacy formData string false "Post privacy"
// @Param location formData string false "Post location"
// @Param media_ids formData int false "Array of mediaIds you want to delete"
// @Param media formData file false "Array of media you want to upload"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{postId} [patch]
func (p *PostUserController) UpdatePost(ctx *gin.Context) {
	var updateInput post_dto.UpdatePostInput

	if err := ctx.ShouldBind(&updateInput); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
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

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, post)
}

// Get many post
// @Summary Get many posts
// @Description Retrieve multiple posts filtered by various criteria.
// @Tags post
// @Accept json
// @Produce json
// @Param userId path string true "User ID to filter posts"
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
// @Router /posts/getMany/{userId} [get]
func (p *PostUserController) GetManyPost(ctx *gin.Context) {
	var query query_object.PostQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	query.UserID = userId

	posts, resultCode, err := services.PostUser().GetManyPost(ctx, &query)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	total := int64(len(posts))

	paging := response.PagingResponse{
		Limit: query.Limit,
		Page:  query.Page,
		Total: total,
	}

	response.SuccessPagingResponse(ctx, response.ErrCodeSuccess, http.StatusOK, posts, paging)
}

// Get post by id documentation
// @Summary Get post by ID
// @Description Retrieve a post by its unique ID
// @Tags post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{postId} [get]
func (p *PostUserController) GetPostById(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	var post *model.Post
	var resultCode int

	post, resultCode, err = services.PostUser().GetPost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, post)
}

// Delete post by id documentation
// @Summary delete post by ID
// @Description when user want to delete post
// @Tags post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/{postId} [delete]
func (p *PostUserController) DeletePost(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	resultCode, err := services.PostUser().DeletePost(ctx, postId)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusNoContent, postId)
}
