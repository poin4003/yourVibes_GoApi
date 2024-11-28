package advertise_user

import (
	"github.com/gin-gonic/gin"
	advertise_services "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	post_services "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdvertise struct {
}

func NewAdvertiseController() *cAdvertise {
	return &cAdvertise{}
}

// CreateAdvertise godoc
// @Summary Comment create advertise
// @Description When user want to create advertise by post
// @Tags advertise_user
// @Accept json
// @Produce json
// @Param input body request.CreateAdvertiseRequest true "input"
// @Security ApiKeyAuth
// @Router /advertise/ [post]
func (c *cAdvertise) CreateAdvertise(ctx *gin.Context) {
	// 1. Get body from form
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to createAdvertiseRequest
	createAdvertiseRequest, ok := body.(*request.CreateAdvertiseRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to check owner
	getOnePostQuery := &query.GetOnePostQuery{
		PostId:              createAdvertiseRequest.PostId,
		AuthenticatedUserId: userIdClaim,
	}

	queryResult, err := post_services.PostUser().GetPost(ctx, getOnePostQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrDataNotFound, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Check owner
	if userIdClaim != queryResult.Post.UserId {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusForbidden, "You can't not promote other people's posts")
		return
	}

	// 6. Call service to handle create advertise
	createAdvertiseCommand, err := createAdvertiseRequest.ToCreateAdvertiseCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := advertise_services.Advertise().CreateAdvertise(ctx, createAdvertiseCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, result.PayUrl)
}
