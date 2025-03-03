package comment_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated request"))
		return
	}

	// 2. Convert to registerRequest
	reportCommentRequest, ok := body.(*request.ReportCommentRequest)
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
	commentReportCommand, err := reportCommentRequest.ToCreateCommentReportCommand(userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}
	result, err := services.CommentReport().CreateCommentReport(ctx, commentReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map result to dto
	commentReportDto := response.ToCommentReportDto(result.CommentReport)

	pkgResponse.OK(ctx, commentReportDto)
}
