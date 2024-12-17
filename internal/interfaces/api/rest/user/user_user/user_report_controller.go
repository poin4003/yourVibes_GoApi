package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
)

type cUserReport struct{}

func NewUserReportController() *cUserReport {
	return &cUserReport{}
}

// ReportUser godoc
// @Summary report user
// @Description When user need to report someone break our rule
// @Tags user_report
// @Accept json
// @Produce json
// @Param input body request.ReportUserRequest true "input"
// @Security ApiKeyAuth
// @Router /users/report [post]
func (c *cUserReport) ReportUser(ctx *gin.Context) {
	// 1. Get body
	body, exists := ctx.Get("validatedRequest")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated request")
		return
	}

	// 2. Convert to registerRequest
	reportUserRequest, ok := body.(*request.ReportUserRequest)
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

	// 4. Call service to handle report user
	userReportCommand, err := reportUserRequest.ToCreateUserReportCommand(userIdClaim)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := services.UserReport().CreateUserReport(ctx, userReportCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map result to dto
	userReportDto := response.ToUserReportDto(result.UserReport)

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, userReportDto)
}
