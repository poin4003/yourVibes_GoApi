package user_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/query"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get reportedUserId from param path
	reportedUserIdStr := ctx.Param("reported_user_id")
	reportedUserId, err := uuid.Parse(reportedUserIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service to handle get user report
	getOneUserReportQuery, err := query.ToGetOneUserReportQuery(userId, reportedUserId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := services.UserReport().GetDetailUserReport(ctx, getOneUserReportQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	userReportDto := response.ToUserReportDto(result.UserReport)

	pkgResponse.OK(ctx, userReportDto)
}

// GetManyUserReports godoc
// @Summary      Get a list of users report
// @Description  Retrieve users report base on filters
// @Tags         admin_user_report
// @Accept       json
// @Produce      json
// @Param        reason              query     string  false  "reason to filter report"
// @Param        status              query     bool    false  "Filter by status"
// @Param        created_at          query     string  false  "Filter by creation day"
// @Param        user_email          query     string  false  "Filter by user email"
// @Param        reported_user_email query     string  false  "Filter by reported user email"
// @Param        admin_email         query     string  false  "Filter by admin email"
// @Param        from_date           query     string  false  "Filter by from date"
// @Param        to_date             query     string  false  "Filter by to date"
// @Param        sort_by       		 query     string  false  "Sort by field"
// @Param        isDescending  		 query     bool    false  "Sort in descending order"
// @Param        limit         		 query     int     false  "Number of results per page"
// @Param        page                query     int     false  "Page number"
// @Security ApiKeyAuth
// @Router       /users/report [get]
func (c *cAdminUserReport) GetManyUserReports(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	userReportQueryObject, ok := queryInput.(*query.UserReportQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3 Call service to handle get many
	getManyUserReportQuery, _ := userReportQueryObject.ToGetManyUserQuery()
	result, err := services.UserReport().GetManyUserReport(ctx, getManyUserReportQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	var userReportDtos []*response.UserReportShortVerDto
	for _, userReportResult := range result.UserReports {
		userReportDtos = append(userReportDtos, response.ToUserReportShortVerDto(userReportResult))
	}

	pkgResponse.OKWithPaging(ctx, userReportDtos, *result.PagingResponse)
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
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get reportedUserId from param path
	reportedUserIdStr := ctx.Param("reported_user_id")
	reportedUserId, err := uuid.Parse(reportedUserIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Get admin id from token
	adminIdClaim, err := extensions.GetAdminID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle user report
	handleUserReportCommand, err := request.ToHandleUserReportCommand(adminIdClaim, userId, reportedUserId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = services.UserReport().HandleUserReport(ctx, handleUserReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. response
	pkgResponse.OK(ctx, nil)
}

// DeleteUserReport godoc
// @Summary delete user report
// @Description When admin need to delete report
// @Tags admin_user_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param reported_user_id path string true "Reported user id"
// @Security ApiKeyAuth
// @Router /users/report/{user_id}/{reported_user_id} [delete]
func (c *cAdminUserReport) DeleteUserReport(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get reportedUserId from param path
	reportedUserIdStr := ctx.Param("reported_user_id")
	reportedUserId, err := uuid.Parse(reportedUserIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service to delete user report
	deleteUserReportCommand, err := request.ToDeleteUserReportCommand(userId, reportedUserId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = services.UserReport().DeleteUserReport(ctx, deleteUserReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. response
	pkgResponse.OK(ctx, nil)
}

// ActivateUserAccount godoc
// @Summary activate user account
// @Description When admin need to activate user account
// @Tags admin_user_report
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Security ApiKeyAuth
// @Router /users/report/activate/{user_id} [patch]
func (c *cAdminUserReport) ActivateUserAccount(ctx *gin.Context) {
	// 1. Get userId from param path
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Call service to activate user account
	activateUserAccountCommand := &command.ActivateUserAccountCommand{
		UserId: userId,
	}

	err = services.UserReport().ActivateUserAccount(ctx, activateUserAccountCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. response
	pkgResponse.OK(ctx, nil)
}
