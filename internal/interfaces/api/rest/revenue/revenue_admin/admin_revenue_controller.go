package revenue_admin

import (
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	revenueServiceQuery "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/query"
	"github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/revenue/revenue_admin/dto/response"
)

type cRevenueAdmin struct{}

func NewRevenueAdminController() *cRevenueAdmin {
	return &cRevenueAdmin{}
}

// GetMonthlyRevenue godoc
// @Summary Get monthly revenue
// @Description Get monthly revenue
// @Tags revenue_admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /revenue/monthly_revenue [get]
func (c *cRevenueAdmin) GetMonthlyRevenue(ctx *gin.Context) {
	// 1. Call service to get monthly revenue
	monthlyRevenueQuery := &revenueServiceQuery.GetMonthlyRevenueQuery{
		Date: time.Now(),
	}

	result, err := services.Revenue().GetMonthlyRevenue(ctx, monthlyRevenueQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Map to dto
	monthlyRevenueDto := response.ToMonthlyRevenueDto(result)

	pkgResponse.OK(ctx, monthlyRevenueDto)
}

// GetSystemStats godoc
// @Summary Get system stats
// @Description Get system stats
// @Tags revenue_admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /revenue/system_stats [get]
func (c *cRevenueAdmin) GetSystemStats(ctx *gin.Context) {
	// 1. Call service to get system stats
	systemStatsQuery := &revenueServiceQuery.GetSystemStatsQuery{
		Date: time.Now(),
	}

	result, err := services.Revenue().GetSystemStats(ctx, systemStatsQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Map to dto
	systemStatsDto := response.ToSystemStatsDto(result)

	pkgResponse.OK(ctx, systemStatsDto)
}
