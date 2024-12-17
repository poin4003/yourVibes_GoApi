package post_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	reportPostRequest, ok := body.(*request.ReportPostRequest)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to handle report post
	postReportCommand, err := reportPostRequest.ToCreatePostReportCommand(userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := services.PostReport().CreatePostReport(ctx, postReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map result to dto
	postReportDto := response.ToPostReportDto(result.PostReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, postReportDto)
}
