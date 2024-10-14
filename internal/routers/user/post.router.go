package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/post_user"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router

	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(authentication.AuthProteced())
	{
		postRouterPrivate.POST("/", post_user.PostUser.CreatePost)
		postRouterPrivate.GET("/getMany/:userId", post_user.PostUser.GetManyPost)
		postRouterPrivate.GET("/:postId", post_user.PostUser.GetPostById)
		postRouterPrivate.PATCH("/:postId", post_user.PostUser.UpdatePost)
		postRouterPrivate.DELETE("/:postId", post_user.PostUser.DeletePost)
	}
}
