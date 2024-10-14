package query_object

import (
	"github.com/google/uuid"
	"time"
)

type UserQueryObject struct {
	UserID       uuid.UUID
	Name         string    `form:"name,omitempty"`
	Email        string    `form:"email,omitempty"`
	PhoneNumber  string    `form:"phone_number,omitempty"`
	Birthday     time.Time `form:"birthday,omitempty"`
	CreatedAt    time.Time `form:"created_at,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"isDescending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}
