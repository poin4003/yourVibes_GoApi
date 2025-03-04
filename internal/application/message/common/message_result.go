package common

import (
	"time"

	"github.com/google/uuid"
)

type MessageResult struct {
	ID        uuid.UUID     `json:"id"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	User      UserResult    `json:"user"`
	Media     []MediaResult `json:"media"`
}

type UserResult struct {
	ID         uuid.UUID `json:"id"`
	FamilyName string    `json:"family_name"`
	Name       string    `json:"name"`
	AvatarUrl  string    `json:"avatar_url"`
}

type MediaResult struct {
	ID        uint      `json:"id"`
	PostId    uuid.UUID `json:"post_id"`
	MediaUrl  string    `json:"media_url"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
