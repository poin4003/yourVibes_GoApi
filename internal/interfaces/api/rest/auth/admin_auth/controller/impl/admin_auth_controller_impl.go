package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/response"
)

type cAdminAuth struct {
	adminAuthService services.IAdminAuth
}

func NewAdminAuthController(
	adminAuthService services.IAdminAuth,
) *cAdminAuth {
	return &cAdminAuth{
		adminAuthService: adminAuthService,
	}
}

// Login admin godoc
// @Summary Admin login
// @Description When user login
// @Tags admin_auth
// @Accept json
// @Produce json
// @Param input body request.AdminLoginRequest true "input"
// @Router /admins/login/ [post]
func (c *cAdminAuth) Login(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to loginRequest
	loginRequest, ok := body.(*request.AdminLoginRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid login request type"))
		return
	}

	// 3. Call service to handle login
	loginCommand, err := loginRequest.ToLoginCommand()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.adminAuthService.Login(ctx, loginCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Convert to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.OK(ctx, gin.H{
		"access_token": result.AccessToken,
		"admin":        adminDto,
	})
}

// ChangeAdminPassword documentation
// @Summary admin change password
// @Description When admin need to change password
// @Tags admin_auth
// @Accept json
// @Produce json
// @Param input body request.ChangeAdminPasswordRequest true "input"
// @Security ApiKeyAuth
// @Router /admins/change_password/ [patch]
func (c *cAdminAuth) ChangeAdminPassword(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to change admin password request
	changeAdminPasswordRequest, ok := body.(*request.ChangeAdminPasswordRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle change password
	changeAdminPasswordCommand, err := changeAdminPasswordRequest.ToChangeAdminPasswordCommand(adminIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = c.adminAuthService.ChangeAdminPassword(ctx, changeAdminPasswordCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}
