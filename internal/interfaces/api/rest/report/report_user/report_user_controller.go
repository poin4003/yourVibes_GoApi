package reportuser

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/services"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user/dto/request"
)

type cReport struct{}

func NewReportController() *cReport {
	return &cReport{}
}

// Report godoc
// @Summary report
// @Description When user need to report break our rule
// @Tags report
// @Accept json
// @Produce json
// @Param input body request.ReportRequest true "input"
// @Security ApiKeyAuth
// @Router /report [post]
func (c *cReport) Report(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(response.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	reportRequest, ok := body.(*request.ReportRequest)
	if !ok {
		ctx.Error(response.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle report post
	reportCommand, err := reportRequest.ToCreateReportCommand(userIdClaim)
	if err != nil {
		ctx.Error(response.NewServerFailedError(err.Error()))
		return
	}
	err = services.Report().CreateReport(ctx, reportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.OK(ctx, nil)
}
