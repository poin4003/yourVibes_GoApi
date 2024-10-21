package post_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cPostShare struct {
}

func NewPostShareController() *cPostShare {
	return &cPostShare{}
}

// SharePost documentation
// @Summary share post
// @Description When user want to share post of another user post's
// @Tags post_share
// @Accept multipart/form-data
// @Produce json
// @Param post_id path string true "PostId"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/share_post/{post_id} [post]
func (p *cPostShare) SharePost(ctx *gin.Context) {
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
		response.ErrorResponse(ctx, response.ErrDataNotFound, http.StatusBadRequest, fmt.Sprint("post id %s not found", postIdStr))
		return
	}

	if userIdClaim == postFound.UserId {
		response.ErrorResponse(ctx, response.ErrDataNotFound, http.StatusBadRequest, fmt.Sprintf("You can not share your own post!"))
		return
	}

	postModel, resultCode, err := services.PostShare().SharePost(ctx, postId, userIdClaim)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, http.StatusInternalServerError, err.Error())
		return
	}

	postDto := mapper.MapPostToNewPostDto(postModel)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postDto)
}
