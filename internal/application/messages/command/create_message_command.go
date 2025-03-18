package command

import (
	"github.com/google/uuid"
	"time"
)

type CreateMessageCommand struct {
	ConversationId uuid.UUID   `json:"conversation_id"`
	UserId         uuid.UUID   `json:"user_id"`
	ParentId       *uuid.UUID  `json:"parent_id"`
	ParentContent  *string     `json:"parent_content"`
	Content        string      `json:"content"`
	User           UserCommand `json:"user"`
	CreatedAt      time.Time   `json:"created_at"`
}

type UserCommand struct {
	ID         string `json:"id"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	AvatarUrl  string `json:"avatar_url"`
}
