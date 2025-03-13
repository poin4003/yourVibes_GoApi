package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	reportAdmin "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin"
	reportAdminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/dto/request"
	reportAdminQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/query"
)

type AdminReportRouter struct{}

func (arr *AdminReportRouter) InitAdminReportRouter(Router *gin.RouterGroup) {
	// Public router
	reportAdminController := reportAdmin.NewAdminReportController()

	// Private router
	reportRouterPrivate := Router.Group("/report")
	reportRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		reportRouterPrivate.GET("/:report_id",
			helpers.ValidateQuery(&reportAdminQuery.ReportDetailQueryObject{}, reportAdminQuery.ValidateReportDetailQueryObject),
			reportAdminController.GetReportDetail,
		)

		reportRouterPrivate.DELETE("/:report_id",
			reportAdminController.DeleteReport,
		)

		reportRouterPrivate.GET("/",
			helpers.ValidateQuery(&reportAdminQuery.ReportQueryObject{}, reportAdminQuery.ValidateReportQueryObject),
			reportAdminController.GetManyReports,
		)

		reportRouterPrivate.PATCH("/handle_report/:report_id",
			helpers.ValidateJsonBody(&reportAdminRequest.HandleReportRequest{}, reportAdminRequest.ValidateHandleReportRequest),
			reportAdminController.HandleReport,
		)

		reportRouterPrivate.PATCH("/activate/:report_id",
			helpers.ValidateJsonBody(&reportAdminRequest.ActivateRequest{}, reportAdminRequest.ValidateActivateRequest),
			reportAdminController.Activate,
		)
	}
}
