package impl

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/application/notification/command"
	pkgResponse "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/query"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	notiService "github.com/poin4003/yourVibes_GoApi/internal/application/notification/services"
	userService "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
)

type cNotification struct {
	notificationService   notiService.INotificationUser
	userInfoService       userService.IUserInfo
	notificationSocketHub *socket_hub.NotificationSocketHub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewNotificationController(
	notificationService notiService.INotificationUser,
	userInfoService userService.IUserInfo,
	notificationSocketHub *socket_hub.NotificationSocketHub,
) *cNotification {
	return &cNotification{
		notificationService:   notificationService,
		userInfoService:       userInfoService,
		notificationSocketHub: notificationSocketHub,
	}
}

// SendNotification documentation
// @Summary Connect to WebSocket
// @Description Establish a WebSocket connection for real-time notifications
// @Tags user_notification
// @Accept json
// @Produce json
// @Router /notification/ws/{user_id} [get]
func (c *cNotification) SendNotification(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(pkgResponse.NewServerFailedError(err.Error()))
		return
	}

	userId := ctx.Param("user_id")
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError(err.Error()))
		conn.Close()
		return
	}

	c.userInfoService.SetUserOnline(ctx, userUUID)
	c.notificationSocketHub.AddConnection(userId, conn)
	fmt.Println("WebSocket connection established")

	go func() {
		defer func() {
			c.userInfoService.SetUserOffline(ctx, userUUID)
			c.notificationSocketHub.RemoveConnection(userId)
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("Unexpected close error:", err)
				} else {
					fmt.Println("WebSocket connection closed for user:", userId, "Error:", err)
				}
				break
			}
			fmt.Printf("Received message from user %s: %s\n", userId, message)
		}
	}()
}

// GetNotification Get notifications
// @Summary Get notifications
// @Tags user_notification
// @Accept json
// @Produce json
// @Param from query string false "Filter notifications by sender"
// @Param notification_type query string false "Filter notifications by type"
// @Param created_at query string false "Filter notifications created at this date"
// @Param sort_by query string false "Sort notifications by this field"
// @Param isDescending query bool false "Sort notifications in descending order"
// @Param limit query int false "Limit the number of notifications returned"
// @Param page query int false "Pagination: page number"
// @Security ApiKeyAuth
// @Router /notification [get]
func (c *cNotification) GetNotification(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(pkgResponse.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	notificationQueryObject, ok := queryInput.(*query.NotificationQueryObject)
	if !ok {
		ctx.Error(pkgResponse.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from param
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle get many
	getManyNotificationQuery, _ := notificationQueryObject.ToGetManyNotificationQuery(userIdClaim)
	result, err := c.notificationService.GetNotificationByUserId(ctx, getManyNotificationQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var notificationDtos []*response.NotificationDto
	for _, notificationResult := range result.Notifications {
		notificationDtos = append(notificationDtos, response.ToNotificationDto(notificationResult))
	}

	pkgResponse.OKWithPaging(ctx, notificationDtos, *result.PagingResponse)
}

// UpdateOneStatusNotifications Update status of notification to false
// @Summary Update notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Security ApiKeyAuth
// @Router /notification/{notification_id} [patch]
func (c *cNotification) UpdateOneStatusNotifications(ctx *gin.Context) {
	// 1. Get notification id from param
	notificationIdStr := ctx.Param("notification_id")
	notificationID, err := strconv.ParseUint(notificationIdStr, 10, 32)
	if err != nil {
		ctx.Error(pkgResponse.NewValidateError("Invalid notification id"))
		return
	}

	// 2. Call service to handle update status notification
	updateOneStatusNotificationCommand := &command.UpdateOneStatusNotificationCommand{
		NotificationId: uint(notificationID),
	}
	err = c.notificationService.UpdateOneStatusNotification(ctx, updateOneStatusNotificationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}

// UpdateManyStatusNotifications Update all status of notification to false
// @Summary Update all notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /notification/ [patch]
func (c *cNotification) UpdateManyStatusNotifications(ctx *gin.Context) {
	// 1. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(pkgResponse.NewInvalidTokenError(err.Error()))
		return
	}

	// 2. Call service to handle update many status by userid
	updateManyStatusNotificationCommand := &command.UpdateManyStatusNotificationCommand{
		UserId: userIdClaim,
	}
	err = c.notificationService.UpdateManyStatusNotification(ctx, updateManyStatusNotificationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	pkgResponse.OK(ctx, nil)
}
