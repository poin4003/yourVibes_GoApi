package user

type UserRouterGroup struct {
	UserRouter
	PostRouter
	CommentRouter
	NotificationRouter
}
