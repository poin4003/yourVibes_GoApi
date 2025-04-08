package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/revenue/revenue_admin/controller"
)

type revenueAdminRouter struct {
	adminRevenueController   controller.IRevenueAdminController
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware
}

func NewRevenueAdminRouter(
	adminRevenueController controller.IRevenueAdminController,
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware,
) *revenueAdminRouter {
	return &revenueAdminRouter{
		adminRevenueController:   adminRevenueController,
		adminProtectedMiddleware: adminProtectedMiddleware,
	}
}

func (r *revenueAdminRouter) InitRevenueAdminRouter(router *gin.RouterGroup) {
	// Private router
	revenueAdminRouterPrivate := router.Group("/revenue")
	revenueAdminRouterPrivate.Use(r.adminProtectedMiddleware.AdminAuthProtected())
	{
		revenueAdminRouterPrivate.GET("/monthly_revenue",
			r.adminRevenueController.GetMonthlyRevenue,
		)

		revenueAdminRouterPrivate.GET("/system_stats",
			r.adminRevenueController.GetSystemStats,
		)
	}
}
