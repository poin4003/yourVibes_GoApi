package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user"
	comment_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
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
		commentRouterPrivate.POST("/", commentUserController.CreateComment)
		commentRouterPrivate.GET("/", commentUserController.GetComment)
		commentRouterPrivate.DELETE("/:comment_id", commentUserController.DeleteComment)
		commentRouterPrivate.PATCH("/:comment_id", commentUserController.UpdateComment)

		// Comment like
		commentRouterPrivate.POST("/like_comment/:comment_id", commentLikeController.LikeComment)
		commentRouterPrivate.GET("/like_comment/:comment_id", commentLikeController.GetUserLikeComment)

		// Comment report
		commentRouterPrivate.POST("/report",
			helpers.ValidateJsonBody(&comment_request.ReportCommentRequest{}, comment_request.ValidateReportCommentRequest),
			commentReportController.ReportComment,
		)
	}
}
