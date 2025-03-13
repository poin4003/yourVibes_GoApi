package advertise_user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	advertiseServices "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	postServiceQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postServices "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/response"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/query"
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
		ctx.Error(response2.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to createAdvertiseRequest
	createAdvertiseRequest, ok := body.(*request.CreateAdvertiseRequest)
	if !ok {
		ctx.Error(response2.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to check owner
	getOnePostQuery := &postServiceQuery.GetOnePostQuery{
		PostId:              createAdvertiseRequest.PostId,
		AuthenticatedUserId: userIdClaim,
	}

	queryResult, err := postServices.PostUser().GetPost(ctx, getOnePostQuery)
	if err != nil {
		ctx.Error(response2.NewDataNotFoundError(err.Error()))
		return
	}

	// 5. Check owner
	if userIdClaim != queryResult.Post.UserId {
		ctx.Error(response2.NewInvalidTokenError("You can't not promote other people's posts"))
		return
	}

	// 6. Check privacy
	if queryResult.Post.Privacy != consts.PUBLIC {
		ctx.Error(response2.NewCustomError(response2.ErrAdMustBePublic, "post privacy is not public"))
		return
	}

	// 7. Call service to handle create advertise
	createAdvertiseCommand, err := createAdvertiseRequest.ToCreateAdvertiseCommand()
	if err != nil {
		ctx.Error(response2.NewServerFailedError(err.Error()))
		return
	}

	result, err := advertiseServices.Advertise().CreateAdvertise(ctx, createAdvertiseCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	response2.OK(ctx, result.PayUrl)
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
		ctx.Error(response2.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to AdvertiseQueryObject
	advertiseQueryObject, ok := queryInput.(*advertiseQuery.AdvertiseQueryObject)
	if !ok {
		ctx.Error(response2.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	postId, err := uuid.Parse(advertiseQueryObject.PostId)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to check owner
	checkPostOwnerQuery := &postServiceQuery.CheckPostOwnerQuery{
		PostId: postId,
		UserId: userIdClaim,
	}
	checkOwnerResult, err := postServices.PostUser().CheckPostOwner(ctx, checkPostOwnerQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	if !checkOwnerResult.IsOwner {
		ctx.Error(response2.NewInvalidTokenError("You can't access this advertise, only for owner"))
		return
	}

	// 5. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseQueryObject.ToGetManyAdvertiseQuery()
	if err != nil {
		ctx.Error(response2.NewServerFailedError(err.Error()))
		return
	}

	result, err := advertiseServices.Advertise().GetManyAdvertise(ctx, getManyAdvertiseQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 6. Convert to dto
	var advertiseDtos []*response.AdvertiseWithBillDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToAdvertiseWithBillDto(*advertiseResult))
	}

	response2.OK(ctx, advertiseDtos)
}
