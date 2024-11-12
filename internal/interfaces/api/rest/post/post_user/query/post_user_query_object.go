package query

import (
	"time"

	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
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

func (req *PostQueryObject) ToGetonePostQuery(
	postId uuid.UUID,
	userId uuid.UUID,
) (*post_query.GetOnePostQuery, error) {
	return &post_query.GetOnePostQuery{
		PostId: postId,
		UserId: userId,
	}, nil
}
