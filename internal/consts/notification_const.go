package consts

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NEW_POST              NotificationType = "new_post"
	NEW_COMMENT           NotificationType = "new_comment"
	LIKE_POST             NotificationType = "like_post"
	LIKE_COMMENT          NotificationType = "like_comment"
	NEW_SHARE             NotificationType = "new_share"
	FRIEND_REQUEST        NotificationType = "friend_request"
	ACCEPT_FRIEND_REQUEST NotificationType = "accept_friend_request"
	DEACTIVATE_POST       NotificationType = "deactivate_post"
	ACTIVATE_POST         NotificationType = "activate_post"
	DEACTIVATE_COMMENT    NotificationType = "deactivate_comment"
	ACTICATE_COMMENT      NotificationType = "activate_comment"
)

var NotificationTypes = []interface{}{
	NEW_POST,
	NEW_COMMENT,
	LIKE_POST,
	LIKE_COMMENT,
	NEW_SHARE,
	FRIEND_REQUEST,
	ACCEPT_FRIEND_REQUEST,
	DEACTIVATE_POST,
	ACTIVATE_POST,
	DEACTIVATE_COMMENT,
	ACTICATE_COMMENT,
}

type NotificationSocketResponse struct {
	ID               uint               `json:"id"`
	From             string             `json:"from"`
	FromUrl          string             `json:"from_url"`
	UserId           uuid.UUID          `json:"user_id"`
	User             UserSocketResponse `json:"user"`
	NotificationType NotificationType   `json:"notification_type"`
	ContentId        string             `json:"content_id"`
	Content          string             `json:"content"`
	Status           bool               `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type UserSocketResponse struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	AvatarUrl  string    `json:"avatar_url"`
}
