package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/controller/notification_controller"
	"github.com/poin4003/yourVibes_GoApi/internal/middlewares/authentication"
)

type NotificationRouter struct{}

func (nr *NotificationRouter) InitNotificationRouter(Router *gin.RouterGroup) {
	notificationController := notification_controller.NewNotificationController()
	// Public router
	notificationRouterPublic := Router.Group("/notifications")
	{
		notificationRouterPublic.GET("/ws/:user_id", notificationController.SendNotification)
	}

	// Private Router
	notificationRouterPrivate := Router.Group("/notifications")
	notificationRouterPrivate.Use(authentication.AuthProteced())
	{
		notificationRouterPrivate.GET("/", notificationController.GetNotification)
		notificationRouterPrivate.PATCH("/:notification_id", notificationController.UpdateOneStatusNotifications)
		notificationRouterPrivate.PATCH("/", notificationController.UpdateManyStatusNotifications)
	}
}
