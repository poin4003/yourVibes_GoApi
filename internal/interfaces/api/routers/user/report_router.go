package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user/controller"
	reportRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user/dto/request"
)

type reportRouter struct {
	reportController        controller.IUserReportController
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware
}

func NewReportRouter(
	reportController controller.IUserReportController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *reportRouter {
	return &reportRouter{
		reportController:        reportController,
		userProtectedMiddleware: userProtectedMiddleware,
	}
}

func (r *reportRouter) InitReportRouter(Router *gin.RouterGroup) {
	// Public router

	// Private router
	reportRouterPrivate := Router.Group("/report")
	reportRouterPrivate.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		reportRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&reportRequest.ReportRequest{}, reportRequest.ValidateReportRequest),
			r.reportController.Report,
		)
	}
}
