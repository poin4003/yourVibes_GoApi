package user

type UserRouterGroup struct {
	userRouter
	postRouter
	commentRouter
	advertiseRouter
	mediaRouter
	messagesRouter
	reportRouter
	notificationRouter
}

func NewUserRouterGroup(
	userRouter userRouter,
	postRouter postRouter,
	commentRouter commentRouter,
	advertiseRouter advertiseRouter,
	mediaRouter mediaRouter,
	messagesRouter messagesRouter,
	reportRouter reportRouter,
	notificationRouter notificationRouter,
) *UserRouterGroup {
	return &UserRouterGroup{
		userRouter:         userRouter,
		postRouter:         postRouter,
		commentRouter:      commentRouter,
		advertiseRouter:    advertiseRouter,
		mediaRouter:        mediaRouter,
		messagesRouter:     messagesRouter,
		reportRouter:       reportRouter,
		notificationRouter: notificationRouter,
	}
}
