package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	cSystemAdmin "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/system/system_admin/controller"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/system/system_admin/dto/request"
)

type systemRouter struct {
	systemAdminCacheController cSystemAdmin.ISystemAdminCache
	systemAdminPostController  cSystemAdmin.ISystemAdminPost
	adminProtectedMiddleware   middlewares.IAdminAuthProtectedMiddleware
}

func NewSystemRouter(
	systemAdminCacheController cSystemAdmin.ISystemAdminCache,
	systemAdminPostController cSystemAdmin.ISystemAdminPost,
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware,
) *systemRouter {
	return &systemRouter{
		systemAdminCacheController: systemAdminCacheController,
		systemAdminPostController:  systemAdminPostController,
		adminProtectedMiddleware:   adminProtectedMiddleware,
	}
}

func (r *systemRouter) InitSystemRouter(Router *gin.RouterGroup) {
	// Public router

	// Private router
	systemRouterPrivate := Router.Group("/systems")
	systemRouterPrivate.Use(r.adminProtectedMiddleware.AdminAuthProtected())
	{
		// System cache router
		systemRouterPrivate.POST("/all_cache",
			r.systemAdminCacheController.ClearAllCache,
		)

		systemRouterPrivate.POST("/post_cache",
			r.systemAdminCacheController.ClearAllPostCache,
		)

		systemRouterPrivate.POST("/comment_cache",
			r.systemAdminCacheController.ClearAllCommentCache,
		)

		// System post router
		systemRouterPrivate.POST("/update_post_and_statistics/:post_id",
			r.systemAdminPostController.UpdatePostAndStatistics,
		)

		systemRouterPrivate.POST("/delay_post_created_at/:post_id",
			r.systemAdminPostController.DelayPostCreatedAt,
		)

		systemRouterPrivate.POST("/expired_advertise/:post_id",
			r.systemAdminPostController.ExpiredAdvertiseByPostId,
		)

		systemRouterPrivate.POST("/push_advertise_to_new_feed",
			helpers.ValidateJsonBody(&request.NumUsers{}, request.ValidateNumUsers),
			r.systemAdminPostController.PushAdvertiseToNewFeed,
		)

		systemRouterPrivate.POST("/push_feature_post_to_new_feed",
			helpers.ValidateJsonBody(&request.NumUsers{}, request.ValidateNumUsers),
			r.systemAdminPostController.PushFeaturePostToNewFeed,
		)

		systemRouterPrivate.POST("/check_expiry_of_advertisement",
			r.systemAdminPostController.CheckExpiryOfAdvertisement,
		)

		systemRouterPrivate.POST("/check_expiry_of_feature_post",
			r.systemAdminPostController.CheckExpiryOfFeaturePost,
		)
	}
}
