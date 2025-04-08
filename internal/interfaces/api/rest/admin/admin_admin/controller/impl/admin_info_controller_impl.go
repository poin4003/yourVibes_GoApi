package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/response"
)

type cAdmin struct {
	adminService services.IAdminInfo
}

func NewAdminController(
	adminService services.IAdminInfo,
) *cAdmin {
	return &cAdmin{
		adminService: adminService,
	}
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	updateAdminInfoRequest, ok := body.(*request.UpdateAdminInfoRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	// 4. Call service to handle update admin
	updateAdminInfoCommand, err := updateAdminInfoRequest.ToUpdateAdminInfoCommand(adminIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.adminService.UpdateAdmin(ctx, updateAdminInfoCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.OK(ctx, adminDto)
}
