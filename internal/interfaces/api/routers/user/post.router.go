package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// Public router
	postUserController := post_user.NewPostUserController(global.Rdb)
	postShareController := post_user.NewPostShareController()
	postLikeController := post_user.NewPostLikeController(global.Rdb)
	postNewFeedController := post_user.NewPostNewFeedController()
	//userRouterPublic := Router.Group("/posts")
	//{
	//}

	// Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(middlewares.AuthProteced())
	{
		// post_user
		postRouterPrivate.POST("/", postUserController.CreatePost)
		postRouterPrivate.GET("/", postUserController.GetManyPost)
		postRouterPrivate.GET("/:post_id", postUserController.GetPostById)
		postRouterPrivate.PATCH("/:post_id", postUserController.UpdatePost)
		postRouterPrivate.DELETE("/:post_id", postUserController.DeletePost)

		// post_like
		postRouterPrivate.POST("/like_post/:post_id", postLikeController.LikePost)
		postRouterPrivate.GET("/like_post/:post_id", postLikeController.GetUserLikePost)

		// post_share
		postRouterPrivate.POST("/share_post/:post_id", postShareController.SharePost)

		// user_new_feed
		postRouterPrivate.DELETE("/new_feeds/:post_id", postNewFeedController.DeleteNewFeed)
		postRouterPrivate.GET("/new_feeds/", postNewFeedController.GetNewFeeds)
	}
}
