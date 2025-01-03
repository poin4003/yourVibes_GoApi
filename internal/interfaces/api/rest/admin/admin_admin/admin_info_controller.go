package admin_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdmin struct{}

func NewAdminController() *cAdmin {
	return &cAdmin{}
}

// UpdateAdminInfo godoc
// @Summary update admin
// @Description When admin need to update info
// @Tags admin
// @Accept json
// @Produce json
// @Param input body request.UpdateAdminInfoRequest true "input"
// @Security ApiKeyAuth
// @Router /admins [patch]
func (c *cAdmin) UpdateAdminInfo(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	updateAdminInfoRequest, ok := body.(*request.UpdateAdminInfoRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Call service to handle update admin
	updateAdminInfoCommand, err := updateAdminInfoRequest.ToUpdateAdminInfoCommand(adminIdClaim)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.AdminInfo().UpdateAdmin(ctx, updateAdminInfoCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, adminDto)
}
