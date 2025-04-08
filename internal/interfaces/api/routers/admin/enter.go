package admin

type AdminRouterGroup struct {
	adminRouter
	advertiseAdminRouter
	revenueAdminRouter
	adminReportRouter
}

func NewAdminRouterGroup(
	adminRouter adminRouter,
	advertiseAdminRouter advertiseAdminRouter,
	revenueAdminRouter revenueAdminRouter,
	adminReportRouter adminReportRouter,
) *AdminRouterGroup {
	return &AdminRouterGroup{
		adminRouter:          adminRouter,
		advertiseAdminRouter: advertiseAdminRouter,
		revenueAdminRouter:   revenueAdminRouter,
		adminReportRouter:    adminReportRouter,
	}
}
