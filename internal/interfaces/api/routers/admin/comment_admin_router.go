package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin"
	admin_comment_report_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin/query"
)

type CommentAdminRouter struct{}

func (car *CommentAdminRouter) InitCommentAdminRouter(Router *gin.RouterGroup) {
	adminCommentReportController := comment_admin.NewAdminCommentReportController()

	// Private router
	adminRouterPrivate := Router.Group("/comments")
	adminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		// comment report
		adminRouterPrivate.GET("/report/:user_id/:reported_comment_id",
			adminCommentReportController.GetCommentReport,
		)

		adminRouterPrivate.GET("/report",
			helpers.ValidateQuery(&admin_comment_report_query.CommentReportQueryObject{}, admin_comment_report_query.ValidateCommentReportQueryObject),
			adminCommentReportController.GetManyCommentReports,
		)
	}
}