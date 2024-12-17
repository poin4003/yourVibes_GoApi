package admin_super_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/query"
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call service to handle get one
	getOneAdminQuery, err := adminRequest.ToGetOneAdminQuery(adminId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := services.SuperAdmin().GetOneAdmin(ctx, getOneAdminQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to Dto
	adminDto := response.ToAdminDto(result.Admin)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, adminDto)
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to AdminQueryObject
	adminQueryObject, ok := queryInput.(*query.AdminQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Call service to handle get many
	getManyAdminQuery, err := adminQueryObject.ToGetManyAdminQuery()

	result, err := services.SuperAdmin().GetManyAdmin(ctx, getManyAdminQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var adminDtos []*response.AdminShortVerResult
	for _, adminResult := range result.Admins {
		adminDtos = append(adminDtos, response.ToAdminShortVerDto(adminResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, http.StatusOK, adminDtos, *result.PagingResponse)
}
