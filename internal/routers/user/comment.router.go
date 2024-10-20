package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/comment_user"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type CommentRouter struct{}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	// Public router

	commentUserController := comment_user.NewCommentUserController()
	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	commentRouterPrivate := Router.Group("/comments")
	commentRouterPrivate.Use(authentication.AuthProteced())
	{
		commentRouterPrivate.POST("/", commentUserController.CreateComment)
		commentRouterPrivate.GET("/", commentUserController.GetComment)
		commentRouterPrivate.DELETE("/:comment_id", commentUserController.DeleteComment)
		commentRouterPrivate.PATCH("/:comment_id", commentUserController.UpdateComment)
	}
}
