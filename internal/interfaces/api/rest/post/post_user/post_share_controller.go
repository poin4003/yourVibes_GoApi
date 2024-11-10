package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
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
// @Param content formData string false "Content of the post"
// @Param privacy formData string false "Privacy level"
// @Param location formData string false "Location of the post"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /posts/share_post/{post_id} [post]
func (p *cPostShare) SharePost(ctx *gin.Context) {
	var sharePostInput request.SharePostInput

	if err := ctx.ShouldBind(&sharePostInput); err != nil {
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

	_, resultCodePostFound, httpStatusCodePostFound, err := services.PostUser().GetPost(ctx, postId, userIdClaim)
	if err != nil {
		response.ErrorResponse(ctx, resultCodePostFound, httpStatusCodePostFound, err.Error())
		return
	}

	postModel, resultCode, httpStatusCode, err := services.PostShare().SharePost(ctx, postId, userIdClaim, &sharePostInput)
	if err != nil {
		response.ErrorResponse(ctx, resultCode, httpStatusCode, err.Error())
		return
	}

	postDto := mapper.MapPostToNewPostDto(postModel)

	response.SuccessResponse(ctx, resultCode, http.StatusOK, postDto)
}
