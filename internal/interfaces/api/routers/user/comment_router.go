package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/controller"
	commentRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/dto/request"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/query"
)

type commentRouter struct {
	commentUserController   controller.ICommentUserController
	commentLikeController   controller.ICommentLikeController
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware
}

func NewCommentRouter(
	commentUserController controller.ICommentUserController,
	commentLikeController controller.ICommentLikeController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *commentRouter {
	return &commentRouter{
		commentUserController:   commentUserController,
		commentLikeController:   commentLikeController,
		userProtectedMiddleware: userProtectedMiddleware,
	}
}

func (r *commentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	// Private router
	commentRouterPrivate := Router.Group("/comments")
	commentRouterPrivate.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		// Comment user
		commentRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&commentRequest.CreateCommentRequest{}, commentRequest.ValidateCreateCommentRequest),
			r.commentUserController.CreateComment,
		)

		commentRouterPrivate.GET("/",
			helpers.ValidateQuery(&commentQuery.CommentQueryObject{}, commentQuery.ValidateCommentQueryObject),
			r.commentUserController.GetComment,
		)

		commentRouterPrivate.PATCH("/:comment_id",
			helpers.ValidateJsonBody(&commentRequest.UpdateCommentRequest{}, commentRequest.ValidateUpdateCommentRequest),
			r.commentUserController.UpdateComment,
		)

		commentRouterPrivate.DELETE("/:comment_id", r.commentUserController.DeleteComment)

		// Comment like
		commentRouterPrivate.POST("/like_comment/:comment_id", r.commentLikeController.LikeComment)

		commentRouterPrivate.GET("/like_comment/:comment_id",
			helpers.ValidateQuery(&commentQuery.CommentLikeQueryObject{}, commentQuery.ValidateCommentLikeQueryObject),
			r.commentLikeController.GetUserLikeComment,
		)
	}
}
