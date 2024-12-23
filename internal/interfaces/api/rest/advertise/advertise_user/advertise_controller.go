package advertise_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	advertise_services "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	post_services "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/response"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cAdvertise struct {
}

func NewAdvertiseController() *cAdvertise {
	return &cAdvertise{}
}

// CreateAdvertise godoc
// @Summary Create advertise
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
	getOnePostQuery := &post_query.GetOnePostQuery{
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

	// 6. Check privacy
	if queryResult.Post.Privacy != consts.PUBLIC {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrAdMustBePublic, http.StatusBadRequest, "post privacy is not public")
		return
	}

	// 7. Call service to handle create advertise
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

// GetManyAdvertise godoc
// @Summary Get many advertise
// @Description Get many advertise
// @Tags advertise_user
// @Accept json
// @Produce json
// @Param post_id query string true "post_id to filter ads"
// @Param limit query int false "Limit of ads per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /advertise/ [get]
func (c *cAdvertise) GetManyAdvertise(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to AdvertiseQueryObject
	advertiseQueryObject, ok := queryInput.(*advertise_query.AdvertiseQueryObject)
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

	postId, err := uuid.Parse(advertiseQueryObject.PostId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to check owner
	checkPostOwnerQuery := &post_query.CheckPostOwnerQuery{
		PostId: postId,
		UserId: userIdClaim,
	}
	checkOwnerResult, err := post_services.PostUser().CheckPostOwner(ctx, checkPostOwnerQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, checkOwnerResult.ResultCode, checkOwnerResult.HttpStatusCode, err.Error())
		return
	}

	if !checkOwnerResult.IsOwner {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusForbidden, "You can't access this advertise, only for owner")
		return
	}

	// 5. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseQueryObject.ToGetManyAdvertiseQuery()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := advertise_services.Advertise().GetManyAdvertise(ctx, getManyAdvertiseQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 6. Convert to dto
	var advertiseDtos []*response.AdvertiseWithBillDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToAdvertiseWithBillDto(*advertiseResult))
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, advertiseDtos)
}
