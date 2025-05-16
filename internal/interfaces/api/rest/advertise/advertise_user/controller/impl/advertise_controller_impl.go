package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	adServices "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	postServiceQuery "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	postServices "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/dto/response"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/query"
)

type cAdvertise struct {
	advertiseServices adServices.IAdvertise
	postService       postServices.IPostUser
}

func NewAdvertiseController(
	advertiseServices adServices.IAdvertise,
	postService postServices.IPostUser,
) *cAdvertise {
	return &cAdvertise{
		advertiseServices: advertiseServices,
		postService:       postService,
	}
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to createAdvertiseRequest
	createAdvertiseRequest, ok := body.(*request.CreateAdvertiseRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to check owner
	getOnePostQuery := &postServiceQuery.GetOnePostQuery{
		PostId:              createAdvertiseRequest.PostId,
		AuthenticatedUserId: userIdClaim,
	}

	queryResult, err := c.postService.GetPost(ctx, getOnePostQuery)
	if err != nil {
		ctx.Error(pkgResponse.NewDataNotFoundError(err.Error()))
		return
	}

	// 5. Check owner
	if userIdClaim != queryResult.Post.UserId {
		ctx.Error(pkgResponse.NewInvalidTokenError("You can't not promote other people's posts"))
		return
	}

	// 6. Check privacy
	if queryResult.Post.Privacy != consts.PUBLIC {
		ctx.Error(pkgResponse.NewCustomError(pkgResponse.ErrAdMustBePublic, "post privacy is not public"))
		return
	}

	// 7. Call service to handle create advertise
	createAdvertiseCommand, err := createAdvertiseRequest.ToCreateAdvertiseCommand()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.advertiseServices.CreateAdvertise(ctx, createAdvertiseCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, result.PayUrl)
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to AdvertiseQueryObject
	advertiseQueryObject, ok := queryInput.(*advertiseQuery.AdvertiseQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	postId, err := uuid.Parse(advertiseQueryObject.PostId)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to check owner
	checkPostOwnerQuery := &postServiceQuery.CheckPostOwnerQuery{
		PostId: postId,
		UserId: userIdClaim,
	}
	checkOwnerResult, err := c.postService.CheckPostOwner(ctx, checkPostOwnerQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	if !checkOwnerResult.IsOwner {
		ctx.Error(pkgResponse.NewInvalidTokenError("You can't access this advertise, only for owner"))
		return
	}

	// 5. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseQueryObject.ToGetManyAdvertiseQuery()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.advertiseServices.GetManyAdvertise(ctx, getManyAdvertiseQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 6. Convert to dto
	var advertiseDtos []*response.AdvertiseWithBillDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToAdvertiseWithBillDto(*advertiseResult))
	}

	pkgResponse.OK(ctx, advertiseDtos)
}

// GetAdvertiseWithStatistic godoc
// @Summary Get advertise with statistic
// @Description Retrieve advertise with statistic
// @Tags advertise_user
// @Accept json
// @Produce json
// @Param advertise_id path string true "Advertise ID"
// @Security ApiKeyAuth
// @Router /advertise/statistic/{advertise_id} [get]
func (c *cAdvertise) GetAdvertiseWithStatistic(ctx *gin.Context) {
	// 1. Get advertise_id from params
	advertiseIdStr := ctx.Param("advertise_id")
	advertiseId, err := uuid.Parse(advertiseIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Call service to get detail
	result, err := c.advertiseServices.GetAdvertiseWithStatistic(ctx, advertiseId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Map to dto
	advertiseDto := response.ToAdvertiseWithStatisticDto(*result)

	pkgResponse.OK(ctx, advertiseDto)
}

// GetAdvertiseByUserId godoc
// @Summary Get many short advertise by user id
// @Description Get many short advertise by user id
// @Tags advertise_user
// @Accept json
// @Produce json
// @Param limit query int false "Limit of ads per page"
// @Param page query int false "Page number for pagination"
// @Security ApiKeyAuth
// @Router /advertise/short_advertise [get]
func (c *cAdvertise) GetAdvertiseByUserId(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to AdvertiseQueryObject
	advertiseByUserIdQueryObject, ok := queryInput.(*advertiseQuery.AdvertiseByUserIdQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 3. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseByUserIdQueryObject.ToGetAdvertiseByUserIdQuery(userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.advertiseServices.GetShortAdvertiseByUserId(ctx, getManyAdvertiseQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 6. Convert to dto
	var advertiseDtos []*response.ShortAdvertiseDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToShortAdvertiseDto(*advertiseResult))
	}

	pkgResponse.OK(ctx, advertiseDtos)
}
