package admin

type AdminRouterGroup struct {
	AdminRouter
	UserAdminRouter
	PostAdminRouter
	CommentAdminRouter
	AdvertiseAdminRouter
	RevenueAdminRouter
}
