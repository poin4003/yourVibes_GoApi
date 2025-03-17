package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/query"
)

type NotificationRouter struct{}

func (nr *NotificationRouter) InitNotificationRouter(Router *gin.RouterGroup) {
	notificationController := notification_user.NewNotificationController()
	// Public router
	notificationRouterPublic := Router.Group("/notification")
	{
		notificationRouterPublic.GET("/ws/:user_id", notificationController.SendNotification)
	}

	// Private router
	notificationRouterPrivate := Router.Group("/notification")
	notificationRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		// notification
		notificationRouterPrivate.GET("/",
			helpers.ValidateQuery(&query.NotificationQueryObject{}, query.ValidateNotificationQueryObject),
			notificationController.GetNotification,
		)

		notificationRouterPrivate.PATCH("/:notification_id", notificationController.UpdateOneStatusNotifications)
		notificationRouterPrivate.PATCH("/", notificationController.UpdateManyStatusNotifications)
	}
}
