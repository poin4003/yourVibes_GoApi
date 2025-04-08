package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/services"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/query"
)

type cAdminReport struct {
	reportService services.IReport
}

func NewAdminReportController(
	reportService services.IReport,
) *cAdminReport {
	return &cAdminReport{
		reportService: reportService,
	}
}

// GetReportDetail documentation
// @Summary Get report detail
// @Description Retrieve a report
// @Tags admin_report
// @Accept json
// @Produce json
// @Param report_id path string true "Report ID"
// @Param report_type query int true "type to get report"
// @Security ApiKeyAuth
// @Router /report/{report_id} [get]
func (c *cAdminReport) GetReportDetail(ctx *gin.Context) {
	// 1. Get reportId from param path
	reportIdStr := ctx.Param("report_id")
	reportId, err := uuid.Parse(reportIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 2. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 3. Convert to reportDetailQueryObject
	reportDetailQueryObject, ok := queryInput.(*query.ReportDetailQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 4. Call service to handle get comment report
	getOneReportQuery, err := reportDetailQueryObject.ToGetOneReportQuery(reportId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	result, err := c.reportService.GetDetailReport(ctx, getOneReportQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	reportDto := response.ToReportResponse(result)

	pkgResponse.OK(ctx, reportDto)
}

// GetManyReports godoc
// @Summary      Get a list of report
// @Description  Retrieve report base on type and filter
// @Tags         admin_report
// @Accept       json
// @Produce      json
// @Param        report_type         query     string  true   "type to get report"
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
// @Router       /report [get]
func (c *cAdminReport) GetManyReports(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	reportQueryObject, ok := queryInput.(*query.ReportQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3 Call service to handle get many
	getManyReportQuery, _ := reportQueryObject.ToGetManyReportQuery()
	result, err := c.reportService.GetManyReport(ctx, getManyReportQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Map to dto
	reportDto := response.ToReportShortVerResponse(result)

	pkgResponse.OKWithPaging(ctx, reportDto, *result.PagingResponse)
}

// HandleReport godoc
// @Summary handle report
// @Description When admin need to handle report
// @Tags admin_report
// @Accept json
// @Produce json
// @Param report_id path string true "Report id"
// @Param input body request.HandleReportRequest true "input"
// @Security ApiKeyAuth
// @Router /report/handle_report/{report_id} [patch]
func (c *cAdminReport) HandleReport(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	handleReportRequest, ok := body.(*request.HandleReportRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 1. Get reportedUserId from param path
	reportIdStr := ctx.Param("report_id")
	reportId, err := uuid.Parse(reportIdStr)
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
	handleUserReportCommand, err := handleReportRequest.ToHandleReportCommand(adminIdClaim, reportId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = c.reportService.HandleReport(ctx, handleUserReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. response
	pkgResponse.OK(ctx, nil)
}

// DeleteReport godoc
// @Summary delete report
// @Description When admin need to delete report
// @Tags admin_report
// @Accept json
// @Produce json
// @Param report_id path string true "Report id"
// @Security ApiKeyAuth
// @Router /report/{report_id} [delete]
func (c *cAdminReport) DeleteReport(ctx *gin.Context) {
	// 1. Get reportedUserId from param path
	reportIdStr := ctx.Param("report_id")
	reportId, err := uuid.Parse(reportIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	// 3. Call service to delete user report
	deleteReportCommand := &command.DeleteReportCommand{ReportId: reportId}
	err = c.reportService.DeleteReport(ctx, deleteReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. response
	pkgResponse.OK(ctx, nil)
}

// Activate godoc
// @Summary activate user account
// @Description When admin need to activate user account
// @Tags admin_report
// @Accept json
// @Produce json
// @Param report_id path string true "report ID"
// @Param input body request.ActivateRequest true "input"
// @Security ApiKeyAuth
// @Router /report/activate/{report_id} [patch]
func (c *cAdminReport) Activate(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	activateRequest, ok := body.(*request.ActivateRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get reportedUserId from param path
	reportIdStr := ctx.Param("report_id")
	reportId, err := uuid.Parse(reportIdStr)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		return
	}

	activateCommand, err := activateRequest.ToActivateCommand(reportId)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	err = c.reportService.Activate(ctx, activateCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. response
	pkgResponse.OK(ctx, nil)
}
