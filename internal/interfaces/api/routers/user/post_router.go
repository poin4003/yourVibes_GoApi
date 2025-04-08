package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/controller"
	postRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/dto/request"
	postQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/query"
)

type postRouter struct {
	postUserController      controller.IPostUserController
	postLikeController      controller.IPostLikeController
	postShareController     controller.IPostShareController
	postNewFeedController   controller.IPostNewFeedController
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware
}

func NewPostRouter(
	postUserController controller.IPostUserController,
	postLikeController controller.IPostLikeController,
	postShareController controller.IPostShareController,
	postNewFeedController controller.IPostNewFeedController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *postRouter {
	return &postRouter{
		postUserController:      postUserController,
		postLikeController:      postLikeController,
		postShareController:     postShareController,
		postNewFeedController:   postNewFeedController,
		userProtectedMiddleware: userProtectedMiddleware,
	}
}

func (r *postRouter) InitPostRouter(Router *gin.RouterGroup) {
	// 1. private router
	postRouterPrivate := Router.Group("/posts")
	postRouterPrivate.Use(r.userProtectedMiddleware.UserAuthProtected())
	{
		// post_user
		postRouterPrivate.POST("/",
			helpers.ValidateFormBody(&postRequest.CreatePostRequest{}, postRequest.ValidateCreatePostRequest),
			r.postUserController.CreatePost,
		)

		postRouterPrivate.GET("/",
			helpers.ValidateQuery(&postQuery.PostQueryObject{}, postQuery.ValidatePostQueryObject),
			r.postUserController.GetManyPost,
		)

		postRouterPrivate.GET("/:post_id", r.postUserController.GetPostById)

		postRouterPrivate.PATCH("/:post_id",
			helpers.ValidateFormBody(&postRequest.UpdatePostRequest{}, postRequest.ValidateUpdatePostRequest),
			r.postUserController.UpdatePost,
		)

		postRouterPrivate.DELETE("/:post_id", r.postUserController.DeletePost)

		// post_like
		postRouterPrivate.POST("/like_post/:post_id", r.postLikeController.LikePost)

		postRouterPrivate.GET("/like_post/:post_id",
			helpers.ValidateQuery(&postQuery.PostLikeQueryObject{}, postQuery.ValidatePostLikeQueryObject),
			r.postLikeController.GetUserLikePost,
		)

		// post_share
		postRouterPrivate.POST("/share_post/:post_id",
			helpers.ValidateFormBody(&postRequest.SharePostRequest{}, postRequest.ValidateSharePostRequest),
			r.postShareController.SharePost,
		)

		// user_new_feed
		postRouterPrivate.DELETE("/new_feeds/:post_id", r.postNewFeedController.DeleteNewFeed)

		postRouterPrivate.GET("/new_feeds/",
			helpers.ValidateQuery(&postQuery.NewFeedQueryObject{}, postQuery.ValidateNewFeedQueryObject),
			r.postNewFeedController.GetNewFeeds,
		)
	}
}
