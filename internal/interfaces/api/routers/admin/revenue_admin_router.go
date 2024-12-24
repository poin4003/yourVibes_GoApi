package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/revenue/revenue_admin"
)

type RevenueAdminRouter struct{}

func (rar *RevenueAdminRouter) InitRevenueAdminRouter(router *gin.RouterGroup) {
	revenueAdminController := revenue_admin.NewRevenueAdminController()

	// Private router
	revenueAdminRouterPrivate := router.Group("/revenue")
	revenueAdminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		revenueAdminRouterPrivate.GET("/monthly_revenue",
			revenueAdminController.GetMonthlyRevenue,
		)

		revenueAdminRouterPrivate.GET("/system_stats",
			revenueAdminController.GetSystemStats,
		)
	}
}
