package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin"
	admin_post_report_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin/query"
)

type PostAdminRouter struct{}

func (par *PostAdminRouter) InitPostAdminRouter(Router *gin.RouterGroup) {
	adminPostReportController := post_admin.NewAdminPostReportController()

	// Private router
	adminRouterPrivate := Router.Group("/posts")
	adminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		// post report
		adminRouterPrivate.GET("/report/:user_id/:reported_post_id",
			adminPostReportController.GetPostReport,
		)

		adminRouterPrivate.GET("/report",
			helpers.ValidateQuery(&admin_post_report_query.PostReportQueryObject{}, admin_post_report_query.ValidatePostReportQueryObject),
			adminPostReportController.GetManyPostReports,
		)

		adminRouterPrivate.PATCH("/report/:user_id/:reported_post_id",
			adminPostReportController.HandlePostReport,
		)
	}
}
