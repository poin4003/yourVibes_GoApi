package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin"
	advertise_admin_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/query"
)

type AdvertiseAdminRouter struct{}

func (aar *AdvertiseAdminRouter) InitAdvertiseAdminRouter(Router *gin.RouterGroup) {
	adminAdvertiseController := advertise_admin.NewAdvertiseAdminController()

	// Private router
	advertiseAdminRouterPrivate := Router.Group("/advertise")
	advertiseAdminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		advertiseAdminRouterPrivate.GET("/:advertise_id",
			adminAdvertiseController.GetAdvertiseDetail,
		)

		advertiseAdminRouterPrivate.GET("/admin",
			helpers.ValidateQuery(&advertise_admin_query.AdvertiseQueryObject{}, advertise_admin_query.ValidateAdvertiseQueryObject),
			adminAdvertiseController.GetManyAdvertise,
		)
	}
}
