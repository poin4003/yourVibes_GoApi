package common

import "github.com/google/uuid"

type PostModerationRequest struct {
	PostID  uuid.UUID `json:"post_id"`
	Content string    `json:"content"`
	BaseURL string    `json:"base_url"`
	Media   []string  `json:"media"`
}

