package user_user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/extensions"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"net/http"
	"strconv"
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
// @Success 200 {string} string "WebSocket connection established"
// @Failure 500 {object} response.ErrResponse
// @Router /users/notifications/ws/{user_id} [get]
func (c *cNotification) SendNotification(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	userId := ctx.Param("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
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
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/notifications [get]
func (c *cNotification) GetNotification(ctx *gin.Context) {
	var query query.NotificationQueryObject

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, err.Error())
		return
	}

	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	getManyNotificationQuery, err := query.ToGetManyNotificationQuery(userIdClaim)

	result, err := services.UserNotification().GetNotificationByUserId(ctx, getManyNotificationQuery)
	if err != nil {
		response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, result.Notifications, *result.PagingResponse)
}

// UpdateOneStatusNotifications Update status of notification to false
// @Summary Update notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/notifications/{notification_id} [patch]
func (c *cNotification) UpdateOneStatusNotifications(ctx *gin.Context) {
	notificationIdStr := ctx.Param("notification_id")
	notificationID, err := strconv.ParseUint(notificationIdStr, 10, 32)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeValidate, http.StatusBadRequest, "Invalid notification id")
		return
	}

	updateOneStatusNotificationCommand := &command.UpdateOneStatusNotificationCommand{
		NotificationId: uint(notificationID),
	}

	result, err := services.UserNotification().UpdateOneStatusNotification(ctx, updateOneStatusNotificationCommand)
	if err != nil {
		response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}

// UpdateManyStatusNotifications Update all status of notification to false
// @Summary Update all notification status to false
// @Tags user_notification
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData
// @Failure 500 {object} response.ErrResponse
// @Security ApiKeyAuth
// @Router /users/notifications/ [patch]
func (c *cNotification) UpdateManyStatusNotifications(ctx *gin.Context) {
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	updatemanyStatusNotificationCommand := &command.UpdateManyStatusNotificationCommand{
		UserId: userIdClaim,
	}

	result, err := services.UserNotification().UpdateManyStatusNotification(ctx, updatemanyStatusNotificationCommand)
	if err != nil {
		response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}
