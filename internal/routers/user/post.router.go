package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/post_user"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router

	postUserController := post_user.NewPostUserController(global.Rdb)
	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(authentication.AuthProteced())
	{
		postRouterPrivate.POST("/", postUserController.CreatePost)
		postRouterPrivate.GET("/getMany/:userId", postUserController.GetManyPost)
		postRouterPrivate.GET("/:postId", postUserController.GetPostById)
		postRouterPrivate.PATCH("/:postId", postUserController.UpdatePost)
		postRouterPrivate.DELETE("/:postId", postUserController.DeletePost)
	}
}
