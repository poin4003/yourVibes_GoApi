package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/controller"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/query"
)

type notificationRouter struct {
	notificationController  controller.IUserNotificationController
	userProtectedController middlewares.IUserAuthProtectedMiddleware
}

func NewNotificationRouter(
	notificationController controller.IUserNotificationController,
	userProtectedMiddleware middlewares.IUserAuthProtectedMiddleware,
) *notificationRouter {
	return &notificationRouter{
		notificationController:  notificationController,
		userProtectedController: userProtectedMiddleware,
	}
}

func (r *notificationRouter) InitNotificationRouter(Router *gin.RouterGroup) {
	// Public router
	notificationRouterPublic := Router.Group("/notification")
	{
		notificationRouterPublic.GET("/ws/:user_id", r.notificationController.SendNotification)
	}

	// Private router
	notificationRouterPrivate := Router.Group("/notification")
	notificationRouterPrivate.Use(r.userProtectedController.UserAuthProtected())
	{
		// notification
		notificationRouterPrivate.GET("/",
			helpers.ValidateQuery(&query.NotificationQueryObject{}, query.ValidateNotificationQueryObject),
			r.notificationController.GetNotification,
		)

		notificationRouterPrivate.PATCH("/:notification_id", r.notificationController.UpdateOneStatusNotifications)
		notificationRouterPrivate.PATCH("/", r.notificationController.UpdateManyStatusNotifications)
	}
}
