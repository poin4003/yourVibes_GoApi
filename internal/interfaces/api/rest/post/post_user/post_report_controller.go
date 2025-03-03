package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type cPostReport struct{}

func NewPostReportController() *cPostReport {
	return &cPostReport{}
}

// ReportPost godoc
// @Summary report post
// @Description When user need to report post break our rule
// @Tags post_report
// @Accept json
// @Produce json
// @Param input body request.ReportPostRequest true "input"
// @Security ApiKeyAuth
// @Router /posts/report [post]
func (c *cPostReport) ReportPost(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	reportPostRequest, ok := body.(*request.ReportPostRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle report post
	postReportCommand, err := reportPostRequest.ToCreatePostReportCommand(userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}
	result, err := services.PostReport().CreatePostReport(ctx, postReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map result to dto
	postReportDto := response.ToPostReportDto(result.PostReport)

	pkgResponse.OK(ctx, postReportDto)
}
