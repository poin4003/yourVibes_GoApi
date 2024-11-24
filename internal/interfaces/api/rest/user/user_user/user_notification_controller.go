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
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/dto/response"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/query"
	pkg_response "github.com/poin4003/yourVibes_GoApi/pkg/response"
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
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, err.Error())
		return
	}

	userId := ctx.Param("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, err.Error())
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
	// 1. Get query
	queryInput, exists := ctx.Get("validatedQuery")
	if !exists {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Missing validated query")
		return
	}

	// 2. Convert to userQueryObject
	notificationQueryObject, ok := queryInput.(*query.NotificationQueryObject)
	if !ok {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrServerFailed, http.StatusInternalServerError, "Invalid register request type")
		return
	}

	// 3. Get user id from param
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Call service to handle get many
	getManyNotificationQuery, err := notificationQueryObject.ToGetManyNotificationQuery(userIdClaim)
	result, err := services.UserNotification().GetNotificationByUserId(ctx, getManyNotificationQuery)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	// 5. Map to dto
	var notificationDtos []*response.NotificationDto
	for _, notificationResult := range result.Notifications {
		notificationDtos = append(notificationDtos, response.ToNotificationDto(notificationResult))
	}

	pkg_response.SuccessPagingResponse(ctx, result.ResultCode, result.HttpStatusCode, notificationDtos, *result.PagingResponse)
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
	// 1. Get notification id from param
	notificationIdStr := ctx.Param("notification_id")
	notificationID, err := strconv.ParseUint(notificationIdStr, 10, 32)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrCodeValidate, http.StatusBadRequest, "Invalid notification id")
		return
	}

	// 2. Call service to handle update status notification
	updateOneStatusNotificationCommand := &command.UpdateOneStatusNotificationCommand{
		NotificationId: uint(notificationID),
	}
	result, err := services.UserNotification().UpdateOneStatusNotification(ctx, updateOneStatusNotificationCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
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
	// 1. Get user id from token
	userIdClaim, err := extensions.GetUserID(ctx)
	if err != nil {
		pkg_response.ErrorResponse(ctx, pkg_response.ErrInvalidToken, http.StatusUnauthorized, err.Error())
		return
	}

	// 2. Call service to handle update many status by userid
	updateManyStatusNotificationCommand := &command.UpdateManyStatusNotificationCommand{
		UserId: userIdClaim,
	}
	result, err := services.UserNotification().UpdateManyStatusNotification(ctx, updateManyStatusNotificationCommand)
	if err != nil {
		pkg_response.ErrorResponse(ctx, result.ResultCode, result.HttpStatusCode, err.Error())
		return
	}

	pkg_response.SuccessResponse(ctx, result.ResultCode, result.HttpStatusCode, nil)
}
