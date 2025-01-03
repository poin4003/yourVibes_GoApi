package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user"
	comment_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
)

type CommentRouter struct{}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	commentUserController := comment_user.NewCommentUserController()
	commentLikeController := comment_user.NewCommentLikeController()
	commentReportController := comment_user.NewCommentReportController()

	// Private router
	commentRouterPrivate := Router.Group("/comments")
	commentRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		// Comment user
		commentRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&comment_request.CreateCommentRequest{}, comment_request.ValidateCreateCommentRequest),
			commentUserController.CreateComment,
		)

		commentRouterPrivate.GET("/",
			helpers.ValidateQuery(&comment_query.CommentQueryObject{}, comment_query.ValidateCommentQueryObject),
			commentUserController.GetComment,
		)

		commentRouterPrivate.PATCH("/:comment_id",
			helpers.ValidateJsonBody(&comment_request.UpdateCommentRequest{}, comment_request.ValidateUpdateCommentRequest),
			commentUserController.UpdateComment,
		)

		commentRouterPrivate.DELETE("/:comment_id", commentUserController.DeleteComment)

		// Comment like
		commentRouterPrivate.POST("/like_comment/:comment_id", commentLikeController.LikeComment)

		commentRouterPrivate.GET("/like_comment/:comment_id",
			helpers.ValidateQuery(&comment_query.CommentLikeQueryObject{}, comment_query.ValidateCommentLikeQueryObject),
			commentLikeController.GetUserLikeComment,
		)

		// Comment report
		commentRouterPrivate.POST("/report",
			helpers.ValidateJsonBody(&comment_request.ReportCommentRequest{}, comment_request.ValidateReportCommentRequest),
			commentReportController.ReportComment,
		)
	}
}
