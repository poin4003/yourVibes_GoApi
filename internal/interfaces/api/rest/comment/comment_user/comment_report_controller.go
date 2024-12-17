package comment_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cCommentReport struct{}

func NewCommentReportController() *cCommentReport {
	return &cCommentReport{}
}

// ReportComment godoc
// @Summary report comment
// @Description When user need to report comment break our rule
// @Tags comment_report
// @Accept json
// @Produce json
// @Param input body request.ReportCommentRequest true "input"
// @Security ApiKeyAuth
// @Router /comments/report [post]
func (c *cCommentReport) ReportComment(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	reportCommentRequest, ok := body.(*request.ReportCommentRequest)
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
	commentReportCommand, err := reportCommentRequest.ToCreateCommentReportCommand(userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := services.CommentReport().CreateCommentReport(ctx, commentReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map result to dto
	commentReportDto := response.ToCommentReportDto(result.CommentReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, commentReportDto)
}
