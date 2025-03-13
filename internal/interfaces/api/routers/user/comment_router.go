package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user"
	commentRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
)

type CommentRouter struct{}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	commentUserController := comment_user.NewCommentUserController()
	commentLikeController := comment_user.NewCommentLikeController()

	// Private router
	commentRouterPrivate := Router.Group("/comments")
	commentRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		// Comment user
		commentRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&commentRequest.CreateCommentRequest{}, commentRequest.ValidateCreateCommentRequest),
			commentUserController.CreateComment,
		)

		commentRouterPrivate.GET("/",
			helpers.ValidateQuery(&commentQuery.CommentQueryObject{}, commentQuery.ValidateCommentQueryObject),
			commentUserController.GetComment,
		)

		commentRouterPrivate.PATCH("/:comment_id",
			helpers.ValidateJsonBody(&commentRequest.UpdateCommentRequest{}, commentRequest.ValidateUpdateCommentRequest),
			commentUserController.UpdateComment,
		)

		commentRouterPrivate.DELETE("/:comment_id", commentUserController.DeleteComment)

		// Comment like
		commentRouterPrivate.POST("/like_comment/:comment_id", commentLikeController.LikeComment)

		commentRouterPrivate.GET("/like_comment/:comment_id",
			helpers.ValidateQuery(&commentQuery.CommentLikeQueryObject{}, commentQuery.ValidateCommentLikeQueryObject),
			commentLikeController.GetUserLikeComment,
		)
	}
}
