package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/controller"
	advertiseAdminQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/query"
)

type advertiseAdminRouter struct {
	advertiseAdminController controller.IAdminAdvertiseController
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware
}

func NewAdvertiseAdminRouter(
	advertiseAdminController controller.IAdminAdvertiseController,
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware,
) *advertiseAdminRouter {
	return &advertiseAdminRouter{
		advertiseAdminController: advertiseAdminController,
		adminProtectedMiddleware: adminProtectedMiddleware,
	}
}

func (r *advertiseAdminRouter) InitAdvertiseAdminRouter(Router *gin.RouterGroup) {
	// Private router
	advertiseAdminRouterPrivate := Router.Group("/advertise")
	advertiseAdminRouterPrivate.Use(r.adminProtectedMiddleware.AdminAuthProtected())
	{
		advertiseAdminRouterPrivate.GET("/:advertise_id",
			r.advertiseAdminController.GetAdvertiseDetail,
		)

		advertiseAdminRouterPrivate.GET("/admin",
			helpers.ValidateQuery(&advertiseAdminQuery.AdvertiseQueryObject{}, advertiseAdminQuery.ValidateAdvertiseQueryObject),
			r.advertiseAdminController.GetManyAdvertise,
		)
	}
}
