package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/post_controller/post_user"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router
	postUserController := post_user.NewPostUserController(global.Rdb)
	postShareController := post_user.NewPostShareController()
	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(authentication.AuthProteced())
	{
		// post_user
		postRouterPrivate.POST("/", postUserController.CreatePost)
		postRouterPrivate.GET("/", postUserController.GetManyPost)
		postRouterPrivate.GET("/:post_id", postUserController.GetPostById)
		postRouterPrivate.PATCH("/:post_id", postUserController.UpdatePost)
		postRouterPrivate.DELETE("/:post_id", postUserController.DeletePost)

		// post_share
		postRouterPrivate.POST("/share_post/:post_id", postShareController.SharePost)
	}
}
