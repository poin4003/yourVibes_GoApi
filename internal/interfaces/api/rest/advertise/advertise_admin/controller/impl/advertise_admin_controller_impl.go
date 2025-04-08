package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	advertiseServiceQuery "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/dto/response"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/query"
)

type cAdvertiseAdmin struct {
	advertiseService services.IAdvertise
}

func NewAdvertiseAdminController(
	advertiseService services.IAdvertise,
) *cAdvertiseAdmin {
	return &cAdvertiseAdmin{
		advertiseService: advertiseService,
	}
}

// GetAdvertiseDetail godoc
// @Summary Get advertise detail
// @Description Retrieve advertise
// @Tags admin_advertise_report
// @Accept json
// @Produce json
// @Param advertise_id path string true "Advertise ID"
// @Security ApiKeyAuth
// @Router /advertise/{advertise_id} [get]
func (c *cAdvertiseAdmin) GetAdvertiseDetail(ctx *gin.Context) {
	// 1. Get advertise_id from params
	advertiseIdStr := ctx.Param("advertise_id")
	advertiseId, err := uuid.Parse(advertiseIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Call service to get detail
	getOneAdvertiseQuery := advertiseServiceQuery.GetOneAdvertiseQuery{
		AdvertiseId: advertiseId,
	}
	result, err := c.advertiseService.GetAdvertise(ctx, &getOneAdvertiseQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Map to dto
	advertiseDto := response.ToAdvertiseDetail(result.Advertise)

	pkgResponse.OK(ctx, advertiseDto)
}

// GetManyAdvertise godoc
// @Summary      Get a list of advertise
// @Description  Retrieve advertise base on filters
// @Tags         admin_advertise_report
// @Accept       json
// @Produce      json
// @Param        post_id             query     string  false  "post id to get advertise"
// @Param        status              query     bool    false  "Filter by status"
// @Param        user_email          query     string  false  "Filter by user email"
// @Param        from_date           query     string  false  "Filter by from date"
// @Param        to_date             query     string  false  "Filter by to date"
// @Param        from_price          query     int     false  "Filter by from price"
// @Param        to_price            query     int     false  "Filter by to price"
// @Param        sort_by       		 query     string  false  "Sort by field"
// @Param        is_descending       query     bool    false  "Sort in descending order"
// @Param        limit         		 query     int     false  "Number of results per page"
// @Param        page                query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /advertise/admin [get]
func (c *cAdvertiseAdmin) GetManyAdvertise(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	advertiseQueryObject, ok := queryInput.(*advertiseQuery.AdvertiseQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseQueryObject.ToGetManyAdvertiseQuery()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	result, err := c.advertiseService.GetManyAdvertise(ctx, getManyAdvertiseQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	var advertiseDtos []*response.AdvertiseWithBillDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToAdvertiseWithBillDto(advertiseResult))
	}

	pkgResponse.OKWithPaging(ctx, advertiseDtos, *result.PagingResponse)
}
