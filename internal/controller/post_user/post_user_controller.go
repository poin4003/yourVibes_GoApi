package post_user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
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

func (p *PostUserController) CreatePost(ctx *gin.Context) {
	var postInput vo.CreatePostInput

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

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, "Authorized")
		return
	}

	userUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	postService := services.PostUser()
	post, resultCode, err := postService.CreatePost(context.Background(), &postInput, uploadedFiles, userUUID)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, http.StatusOK, post)
}
