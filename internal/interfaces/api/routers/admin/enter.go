package admin

type AdminRouterGroup struct {
	adminRouter
	advertiseAdminRouter
	revenueAdminRouter
	adminReportRouter
	systemRouter
}

func NewAdminRouterGroup(
	adminRouter adminRouter,
	advertiseAdminRouter advertiseAdminRouter,
	revenueAdminRouter revenueAdminRouter,
	adminReportRouter adminReportRouter,
	systemRouter systemRouter,
) *AdminRouterGroup {
	return &AdminRouterGroup{
		adminRouter:          adminRouter,
		advertiseAdminRouter: advertiseAdminRouter,
		revenueAdminRouter:   revenueAdminRouter,
		adminReportRouter:    adminReportRouter,
		systemRouter:         systemRouter,
	}
}
