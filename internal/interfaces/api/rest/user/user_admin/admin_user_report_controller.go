package user_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cAdminUserReport struct{}

func NewAdminUserReportController() *cAdminUserReport {
	return &cAdminUserReport{}
}

// GetUserReport documentation
// @Summary Get user report detail
// @Description Retrieve a user report
// @Tags admin_user_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_user_id path string true "Reported user id"
// @Security ApiKeyAuth
// @Router /users/report/{user_id}/{reported_user_id} [get]
func (c *cAdminUserReport) GetUserReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedUserId from param path
	reportedUserIdStr := ctx.Param("reported_user_id")
	reportedUserId, err := uuid.Parse(reportedUserIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call service to handle get user report
	getOneUserReportQuery, err := query.ToGetOneUserReportQuery(userId, reportedUserId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserReport().GetDetailUserReport(ctx, getOneUserReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	userReportDto := response.ToUserReportDto(result.UserReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, http.StatusOK, userReportDto)
}

// GetManyUserReports godoc
// @Summary      Get a list of users report
// @Description  Retrieve users report base on filters
// @Tags         admin_user_report
// @Accept       json
// @Produce      json
// @Param        reason        query     string  false  "reason to filter report"
// @Param        status        query     bool    false  "Filter by status"
// @Param        created_at    query     string  false  "Filter by creation day"
// @Param        sort_by       query     string  false  "Sort by field"
// @Param        isDescending  query     bool    false  "Sort in descending order"
// @Param        limit         query     int     false  "Number of results per page"
// @Param        page          query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /users/report [get]
func (c *cAdminUserReport) GetManyUserReports(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to userQueryObject
	userReportQueryObject, ok := queryInput.(*query.UserReportQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3 Call service to handle get many
	getManyUserReportQuery, err := userReportQueryObject.ToGetManyUserQuery()
	result, err := services.UserReport().GetManyUserReport(ctx, getManyUserReportQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. Map to dto
	var userReportDtos []*response.UserReportShortVerDto
	for _, userReportResult := range result.UserReports {
		userReportDtos = append(userReportDtos, response.ToUserReportShortVerDto(userReportResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, userReportDtos, *result.PagingResponse)
}

// HandleUserReport godoc
// @Summary handle user report
// @Description When admin need to handle report
// @Tags admin_user_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_user_id path string true "Reported user id"
// @Security ApiKeyAuth
// @Router /users/report/{user_id}/{reported_user_id} [patch]
func (c *cAdminUserReport) HandleUserReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Get reportedUserId from param path
	reportedUserIdStr := ctx.Param("reported_user_id")
	reportedUserId, err := uuid.Parse(reportedUserIdStr)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Call service to handle user report
	handleUserReportCommand, err := request.ToHandleUserReportCommand(adminIdClaim, userId, reportedUserId)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := services.UserReport().HandleUserReport(ctx, handleUserReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 4. response
	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}
