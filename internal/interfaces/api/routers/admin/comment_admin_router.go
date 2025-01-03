package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin"
	adminCommentReportQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_admin/query"
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
			helpers.ValidateQuery(&adminCommentReportQuery.CommentReportQueryObject{}, adminCommentReportQuery.ValidateCommentReportQueryObject),
			adminCommentReportController.GetManyCommentReports,
		)

		adminRouterPrivate.PATCH("/report/:user_id/:reported_comment_id",
			adminCommentReportController.HandleCommentReport,
		)

		adminRouterPrivate.DELETE("/report/:user_id/:reported_comment_id",
			adminCommentReportController.DeleteCommentReport,
		)

		adminRouterPrivate.PATCH("/report/activate/:comment_id",
			adminCommentReportController.ActivateComment,
		)
	}
}
