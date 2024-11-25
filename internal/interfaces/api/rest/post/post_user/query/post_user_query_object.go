package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
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

func ValidatePostQueryObject(input interface{}) error {
	query, ok := input.(*PostQueryObject)
	if !ok {
		return fmt.Errorf("validate PostQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Content, validation.Min(1)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *PostQueryObject) ToGetOnePostQuery(
	postId uuid.UUID,
	authenticatedUserId uuid.UUID,
) (*post_query.GetOnePostQuery, error) {
	return &post_query.GetOnePostQuery{
		PostId:              postId,
		AuthenticatedUserId: authenticatedUserId,
	}, nil
}

func (req *PostQueryObject) ToGetManyPostQuery(
	authenticatedUserId uuid.UUID,
) (*post_query.GetManyPostQuery, error) {
	var userId uuid.UUID
	if req.UserID != "" {
		parseUserId, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, err
		}
		userId = parseUserId
	}

	return &post_query.GetManyPostQuery{
		AuthenticatedUserId: authenticatedUserId,
		UserID:              userId,
		Content:             req.Content,
		Location:            req.Location,
		IsAdvertisement:     req.IsAdvertisement,
		CreatedAt:           req.CreatedAt,
		SortBy:              req.SortBy,
		IsDescending:        req.IsDescending,
		Limit:               req.Limit,
		Page:                req.Page,
	}, nil
}
