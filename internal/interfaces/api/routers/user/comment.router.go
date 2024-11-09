package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/controller/comment/comment_user"
)

type CommentRouter struct{}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	// Public router

	commentUserController := comment_user.NewCommentUserController()
	commentLikeController := comment_user.NewCommentLikeController()
	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	commentRouterPrivate := Router.Group("/comments")
	commentRouterPrivate.Use(middlewares.AuthProteced())
	{
		// Comment user
		commentRouterPrivate.POST("/", commentUserController.CreateComment)
		commentRouterPrivate.GET("/", commentUserController.GetComment)
		commentRouterPrivate.DELETE("/:comment_id", commentUserController.DeleteComment)
		commentRouterPrivate.PATCH("/:comment_id", commentUserController.UpdateComment)

		// Comment like
		commentRouterPrivate.POST("/like_comment/:comment_id", commentLikeController.LikeComment)
		commentRouterPrivate.GET("/like_comment/:comment_id", commentLikeController.GetUserLikeComment)
	}
}
