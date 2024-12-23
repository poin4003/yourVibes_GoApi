package advertise_admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	query_service "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/dto/response"
	query_interface "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cAdvertiseAdmin struct{}

func NewAdvertiseAdminController() *cAdvertiseAdmin {
	return &cAdvertiseAdmin{}
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Call service to get detail
	getOneAdvertiseQuery := query_service.GetOneAdvertiseQuery{
		AdvertiseId: advertiseId,
	}
	result, err := services.Advertise().GetAdvertise(ctx, &getOneAdvertiseQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 3. Map to dto
	advertiseDto := response.ToAdvertiseDetail(result.Advertise)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, advertiseDto)
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to userQueryObject
	advertiseQueryObject, ok := queryInput.(*query_interface.AdvertiseQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle get many
	getManyAdvertiseQuery, err := advertiseQueryObject.ToGetManyAdvertiseQuery()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	result, err := services.Advertise().GetManyAdvertise(ctx, getManyAdvertiseQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var advertiseDtos []*response.AdvertiseWithBillDto
	for _, advertiseResult := range result.Advertises {
		advertiseDtos = append(advertiseDtos, response.ToAdvertiseWithBillDto(advertiseResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, advertiseDtos, *result.PagingResponse)
}
