package user

type UserRouterGroup struct {
	UserRouter
	PostRouter
	CommentRouter
	AdvertiseRouter
	MediaRouter
	MessagesRouter
	ReportRouter
	NotificationRouter
}
