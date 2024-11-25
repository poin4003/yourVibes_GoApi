package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user"
	post_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
)

type PostRouter struct{}

func (pr *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	// 1. Init Controller
	postUserController := post_user.NewPostUserController(global.Rdb)
	postShareController := post_user.NewPostShareController()
	postLikeController := post_user.NewPostLikeController(global.Rdb)
	postNewFeedController := post_user.NewPostNewFeedController()

	// 2. Private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(middlewares.AuthProteced())
	{
		// post_user
		postRouterPrivate.POST("/",
			helpers.ValidateFormBody(&post_request.CreatePostRequest{}, post_request.ValidateCreatePostRequest),
			postUserController.CreatePost,
		)

		postRouterPrivate.GET("/",
			helpers.ValidateQuery(&post_query.PostQueryObject{}, post_query.ValidatePostQueryObject),
			postUserController.GetManyPost,
		)

		postRouterPrivate.GET("/:post_id", postUserController.GetPostById)

		postRouterPrivate.PATCH("/:post_id",
			helpers.ValidateFormBody(&post_request.UpdatePostRequest{}, post_request.ValidateUpdatePostRequest),
			postUserController.UpdatePost,
		)

		postRouterPrivate.DELETE("/:post_id", postUserController.DeletePost)

		// post_like
		postRouterPrivate.POST("/like_post/:post_id", postLikeController.LikePost)

		postRouterPrivate.GET("/like_post/:post_id",
			helpers.ValidateQuery(&post_query.PostLikeQueryObject{}, post_query.ValidatePostLikeQueryObject),
			postLikeController.GetUserLikePost,
		)

		// post_share
		postRouterPrivate.POST("/share_post/:post_id", postShareController.SharePost)

		// user_new_feed
		postRouterPrivate.DELETE("/new_feeds/:post_id", postNewFeedController.DeleteNewFeed)

		postRouterPrivate.GET("/new_feeds/",
			helpers.ValidateQuery(&post_query.NewFeedQueryObject{}, post_query.ValidateNewFeedQueryObject),
			postNewFeedController.GetNewFeeds,
		)
	}
}
