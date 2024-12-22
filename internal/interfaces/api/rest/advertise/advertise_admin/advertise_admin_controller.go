package advertise_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdvertiseAdmin struct{}

func NewAdvertiseAdminController() *cAdvertiseAdmin {
	return &cAdvertiseAdmin{}
}

func (c *cAdvertiseAdmin) GetAdvertiseDetail(ctx *gin.Context) {
	// 1. Get advertise_id from params
	advertiseIdStr := ctx.Param("advertise_id")
	advertiseId, err := uuid.Parse(advertiseIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Call service to get detail
	getOneAdvertiseQuery := query.GetOneAdvertiseQuery{
		AdvertiseId: advertiseId,
	}
	result, err := services.Advertise().GetAdvertise(ctx, &getOneAdvertiseQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 3. Map to dto
	//advertiseDto :=
}
