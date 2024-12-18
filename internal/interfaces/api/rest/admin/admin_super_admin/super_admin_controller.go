package admin_super_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cSuperAdmin struct{}

func NewSuperAdminController() *cSuperAdmin {
	return &cSuperAdmin{}
}

// CreateAdmin godoc
// @Summary Create admin
// @Description When super admin need to create new admin
// @Tags super_admin
// @Accept json
// @Produce json
// @Param input body request.CreateAdminRequest true "input"
// @Security ApiKeyAuth
// @Router /admins/super_admin [post]
func (c *cSuperAdmin) CreateAdmin(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	createAdminRequest, ok := body.(*request.CreateAdminRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle create admin
	createAdminCommand, err := createAdminRequest.ToCreateAdminCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.SuperAdmin().CreateAdmin(ctx, createAdminCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, adminDto)
}

// UpdateAdmin godoc
// @Summary update admin
// @Description When super admin need to update role and status of admin
// @Tags super_admin
// @Accept json
// @Produce json
// @Param input body request.UpdateAdminForSuperAdminRequest true "input"
// @Security ApiKeyAuth
// @Router /admins/super_admin [patch]
func (c *cSuperAdmin) UpdateAdmin(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	updateAdminForSuperAdminRequest, ok := body.(*request.UpdateAdminForSuperAdminRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle update admin
	updateAdminForSuperAdminCommand, err := updateAdminForSuperAdminRequest.ToUpdateAdminForSuperAdminCommand()
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.SuperAdmin().UpdateAdmin(ctx, updateAdminForSuperAdminCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, adminDto)
}
