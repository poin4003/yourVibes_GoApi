package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
)

type cPostShare struct {
	postShareService services.IPostShare
}

func NewPostShareController(
	postShareService services.IPostShare,
) *cPostShare {
	return &cPostShare{
		postShareService: postShareService,
	}
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
// @Security ApiKeyAuth
// @Router /posts/share_post/{post_id} [post]
func (c *cPostShare) SharePost(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to updateUserRequest
	sharePostRequest, ok := body.(*request.SharePostRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 2. Get post id from param
	postIdStr := ctx.Param("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle sharing
	sharePostCommand, err := sharePostRequest.ToSharePostCommand(postId, userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.postShareService.SharePost(ctx, sharePostCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	postDto := response.ToPostDto(*result.Post)

	pkgResponse.OK(ctx, postDto)
}
