package admin_super_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	createAdminRequest, ok := body.(*request.CreateAdminRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Call service to handle create admin
	createAdminCommand, err := createAdminRequest.ToCreateAdminCommand()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.SuperAdmin().CreateAdmin(ctx, createAdminCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.OK(ctx, adminDto)
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	updateAdminForSuperAdminRequest, ok := body.(*request.UpdateAdminForSuperAdminRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Call service to handle update admin
	updateAdminForSuperAdminCommand, err := updateAdminForSuperAdminRequest.ToUpdateAdminForSuperAdminCommand()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.SuperAdmin().UpdateAdmin(ctx, updateAdminForSuperAdminCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map result to dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.OK(ctx, adminDto)
}

// GetAdminById documentation
// @Summary Get admin by ID
// @Description Retrieve admin by its unique ID
// @Tags super_admin
// @Accept json
// @Produce json
// @Param admin_id path string true "Admin ID"
// @Security ApiKeyAuth
// @Router /admins/{admin_id} [get]
func (c *cSuperAdmin) GetAdminById(ctx *gin.Context) {
	var adminRequest query.AdminQueryObject

	// 1. Get post id from param
	adminIdStr := ctx.Param("admin_id")
	adminId, err := uuid.Parse(adminIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service to handle get one
	getOneAdminQuery, err := adminRequest.ToGetOneAdminQuery(adminId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}
	result, err := services.SuperAdmin().GetOneAdmin(ctx, getOneAdminQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to Dto
	adminDto := response.ToAdminDto(result.Admin)

	pkgResponse.OK(ctx, adminDto)
}

// GetManyAdmins godoc
// @Summary      Get a list of admins
// @Description  Retrieve admins based on filters
// @Tags         super_admin
// @Accept       json
// @Produce      json
// @Param        name          query     string  false  "name to filter admins"
// @Param        email         query     string  false  "Filter by email"
// @Param        phone_number  query     string  false  "Filter by phone number"
// @Param        identity_id   query     string  false  "Filter by identity id"
// @Param        birthday      query     string  false  "Filter by birthday"
// @Param        created_at    query     string  false  "Filter by creation day"
// @Param        status        query     bool    false  "Filter by status"
// @Param        role          query     bool    false  "Filter by role"
// @Param        sort_by       query     string  false  "Sort by field"
// @Param        is_descending  query     bool    false  "Sort in descending order"
// @Param        limit         query     int     false  "Number of results per page"
// @Param        page          query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /admins/ [get]
func (c *cSuperAdmin) GetManyAdmins(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to AdminQueryObject
	adminQueryObject, ok := queryInput.(*query.AdminQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Call service to handle get many
	getManyAdminQuery, _ := adminQueryObject.ToGetManyAdminQuery()

	result, err := services.SuperAdmin().GetManyAdmin(ctx, getManyAdminQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	var adminDtos []*response.AdminDto
	for _, adminResult := range result.Admins {
		adminDtos = append(adminDtos, response.ToAdminDto(adminResult))
	}

	pkgResponse.OKWithPaging(ctx, adminDtos, *result.PagingResponse)
}

// ForgotAdminPassword documentation
// @Summary Admin forgot password
// @Description When super admin change admin password
// @Tags super_admin
// @Accept json
// @Produce json
// @Param input body request.ForgotAdminPasswordRequest true "input"
// @Security ApiKeyAuth
// @Router /admins/super_admin/forgot_admin_password [post]
func (c *cSuperAdmin) ForgotAdminPassword(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to forgot admin password request
	forgotAdminPasswordRequest, ok := body.(*request.ForgotAdminPasswordRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Call service to handle forgot  password
	forgotAdminPasswordCommand, err := forgotAdminPasswordRequest.ToForgotAdminPasswordCommand()
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = services.AdminAuth().ForgotAdminPassword(ctx, forgotAdminPasswordCommand)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	pkgResponse.OK(ctx, nil)
}
