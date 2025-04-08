package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/controller"
	reportAdminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/dto/request"
	reportAdminQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/query"
)

type adminReportRouter struct {
	adminReportController    controller.IAdminReportController
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware
}

func NewAdminReportRouter(
	adminReportController controller.IAdminReportController,
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware,
) *adminReportRouter {
	return &adminReportRouter{
		adminReportController:    adminReportController,
		adminProtectedMiddleware: adminProtectedMiddleware,
	}
}

func (r *adminReportRouter) InitAdminReportRouter(Router *gin.RouterGroup) {
	// Public router
	// Private router
	reportRouterPrivate := Router.Group("/report")
	reportRouterPrivate.Use(r.adminProtectedMiddleware.AdminAuthProtected())
	{
		reportRouterPrivate.GET("/:report_id",
			helpers.ValidateQuery(&reportAdminQuery.ReportDetailQueryObject{}, reportAdminQuery.ValidateReportDetailQueryObject),
			r.adminReportController.GetReportDetail,
		)

		reportRouterPrivate.DELETE("/:report_id",
			r.adminReportController.DeleteReport,
		)

		reportRouterPrivate.GET("/",
			helpers.ValidateQuery(&reportAdminQuery.ReportQueryObject{}, reportAdminQuery.ValidateReportQueryObject),
			r.adminReportController.GetManyReports,
		)

		reportRouterPrivate.PATCH("/handle_report/:report_id",
			helpers.ValidateJsonBody(&reportAdminRequest.HandleReportRequest{}, reportAdminRequest.ValidateHandleReportRequest),
			r.adminReportController.HandleReport,
		)

		reportRouterPrivate.PATCH("/activate/:report_id",
			helpers.ValidateJsonBody(&reportAdminRequest.ActivateRequest{}, reportAdminRequest.ValidateActivateRequest),
			r.adminReportController.Activate,
		)
	}
}
