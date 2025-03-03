package user_user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
		ctx.Error(pkgResponse.NewServerFailedError("Missing validate request"))
		return
	}

	// 2. Convert to registerRequest
	reportUserRequest, ok := body.(*request.ReportUserRequest)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid request type"))
		return
	}

	// 3. Get userId from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle report user
	userReportCommand, err := reportUserRequest.ToCreateUserReportCommand(userIdClaim)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError())
		return
	}
	result, err := services.UserReport().CreateUserReport(ctx, userReportCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map result to dto
	userReportDto := response.ToUserReportDto(result.UserReport)

	pkgResponse.OK(ctx, userReportDto)
}
