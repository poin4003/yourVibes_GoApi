package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin"
	adminPostReportQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_admin/query"
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
			helpers.ValidateQuery(&adminPostReportQuery.PostReportQueryObject{}, adminPostReportQuery.ValidatePostReportQueryObject),
			adminPostReportController.GetManyPostReports,
		)

		adminRouterPrivate.PATCH("/report/:user_id/:reported_post_id",
			adminPostReportController.HandlePostReport,
		)

		adminRouterPrivate.DELETE("/report/:user_id/:reported_post_id",
			adminPostReportController.DeletePostReport,
		)

		adminRouterPrivate.PATCH("/report/activate/:post_id",
			adminPostReportController.ActivatePost,
		)
	}
}
