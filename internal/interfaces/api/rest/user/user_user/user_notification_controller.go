package user_user

import (
	"fmt"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
)

type cNotification struct{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewNotificationController() *cNotification {
	return &cNotification{}
}

// SendNotification documentation
// @Summary Connect to WebSocket
// @Description Establish a WebSocket connection for real-time notifications
// @Tags user_notification
// @Accept json
// @Produce json
// @Router /users/notifications/ws/{user_id} [get]
func (c *cNotification) SendNotification(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(response2.NewServerFailedError(err.Error()))
		return
	}

	userId := ctx.Param("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		ctx.Error(response2.NewValidateError(err.Error()))
		conn.Close()
		return
	}

	global.SocketHub.AddConnection(userId, conn)
	fmt.Println("WebSocket connection established")

	go func() {
		defer global.SocketHub.RemoveConnection(userId)
		defer conn.Close()

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
// @Router /users/notifications [get]
func (c *cNotification) GetNotification(ctx *gin.Context) {
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		ctx.Error(response2.NewServerFailedError("Missing validated query"))
		return
	}

	// 2. Convert to userQueryObject
	notificationQueryObject, ok := queryInput.(*query.NotificationQueryObject)
	if !ok {
		ctx.Error(response2.NewServerFailedError("Invalid register request type"))
		return
	}

	// 3. Get user id from param
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	// 4. Call service to handle get many
	getManyNotificationQuery, _ := notificationQueryObject.ToGetManyNotificationQuery(userIdClaim)
	result, err := services.UserNotification().GetNotificationByUserId(ctx, getManyNotificationQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Map to dto
	var notificationDtos []*response.NotificationDto
	for _, notificationResult := range result.Notifications {
		notificationDtos = append(notificationDtos, response.ToNotificationDto(notificationResult))
	}

	response2.OKWithPaging(ctx, notificationDtos, *result.PagingResponse)
}

// UpdateOneStatusNotifications Update status of notification to false
// @Summary Update notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Security ApiKeyAuth
// @Router /users/notifications/{notification_id} [patch]
func (c *cNotification) UpdateOneStatusNotifications(ctx *gin.Context) {
	// 1. Get notification id from param
	notificationIdStr := ctx.Param("notification_id")
	notificationID, err := strconv.ParseUint(notificationIdStr, 10, 32)
	if err != nil {
		ctx.Error(response2.NewValidateError("Invalid notification id"))
		return
	}

	// 2. Call service to handle update status notification
	updateOneStatusNotificationCommand := &command.UpdateOneStatusNotificationCommand{
		NotificationId: uint(notificationID),
	}
	err = services.UserNotification().UpdateOneStatusNotification(ctx, updateOneStatusNotificationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	response2.OK(ctx, nil)
}

// UpdateManyStatusNotifications Update all status of notification to false
// @Summary Update all notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /users/notifications/ [patch]
func (c *cNotification) UpdateManyStatusNotifications(ctx *gin.Context) {
	// 1. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		ctx.Error(response2.NewInvalidTokenError(err.Error()))
		return
	}

	// 2. Call service to handle update many status by userid
	updateManyStatusNotificationCommand := &command.UpdateManyStatusNotificationCommand{
		UserId: userIdClaim,
	}
	err = services.UserNotification().UpdateManyStatusNotification(ctx, updateManyStatusNotificationCommand)
	if err != nil {
		ctx.Error(err)
		return
	}

	response2.OK(ctx, nil)
}
