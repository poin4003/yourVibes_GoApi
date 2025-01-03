package admin_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdminAuth struct{}

func NewAdminAuthController() *cAdminAuth {
	return &cAdminAuth{}
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
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to loginRequest
	loginRequest, ok := body.(*request.AdminLoginRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid login request type")
		return
	}

	// 3. Call service to handle login
	loginCommand, err := loginRequest.ToLoginCommand()
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.AdminAuth().Login(ctx, loginCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Convert to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, gin.H{
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
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to change admin password request
	changeAdminPasswordRequest, ok := body.(*request.ChangeAdminPasswordRequest)
	if !ok {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to handle change password
	changeAdminPasswordCommand, err := changeAdminPasswordRequest.ToChangeAdminPasswordCommand(adminIdClaim)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, pkgResponse.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.AdminAuth().ChangeAdminPassword(ctx, changeAdminPasswordCommand)
	if err != nil {
		pkgResponse.ErrorResponse(ctx, result.ResultCode, http.StatusBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessResponse(ctx, pkgResponse.ErrCodeSuccess, http.StatusOK, nil)
}
