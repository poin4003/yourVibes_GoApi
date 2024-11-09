package query

import (
	"time"
)

type PostQueryObject struct {
	UserID          string    `form:"user_id,omitempty"`
	Content         string    `form:"content,omitempty"`
	Location        string    `form:"location,omitempty"`
	IsAdvertisement bool      `form:"is_advertisement,omitempty"`
	CreatedAt       time.Time `form:"created_at,omitempty"`
	SortBy          string    `form:"sort_by,omitempty"`
	IsDescending    bool      `form:"isDescending,omitempty"`
	Limit           int       `form:"limit,omitempty"`
	Page            int       `form:"page,omitempty"`
}
