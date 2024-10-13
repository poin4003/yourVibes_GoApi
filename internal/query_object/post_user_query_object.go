package query_object

import (
	"github.com/google/uuid"
	"time"
)

type PostQueryObject struct {
	UserID          uuid.UUID
	Title           string    `form:"title,omitempty"`
	Content         string    `form:"content,omitempty"`
	Location        string    `form:"location,omitempty"`
	IsAdvertisement bool      `form:"is_advertisement,omitempty"`
	CreatedAt       time.Time `form:"created_at,omitempty"`
	SortBy          string    `form:"sort_by,omitempty"`
	IsDescending    bool      `form:"isDescending,omitempty"`
	Limit           int       `form:"limit,omitempty"`
	Page            int       `form:"page,omitempty"`
}
