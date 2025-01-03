package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
)

type AdvertiseQueryObject struct {
	PostId       string    `form:"post_id,omitempty"`
	UserEmail    string    `form:"user_email,omitempty"`
	Status       *bool     `form:"status,omitempty"`
	FromDate     time.Time `form:"from_date,omitempty"`
	ToDate       time.Time `form:"to_date,omitempty"`
	FromPrice    int       `form:"from_price,omitempty"`
	ToPrice      int       `form:"to_price,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"is_descending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}

func ValidateAdvertiseQueryObject(input interface{}) error {
	query, ok := input.(*AdvertiseQueryObject)
	if !ok {
		return fmt.Errorf("validateAdvertiseQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.UserEmail, is.Email),
		validation.Field(&query.FromPrice, validation.Min(0)),
		validation.Field(&query.ToPrice, validation.Min(0)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *AdvertiseQueryObject) ToGetManyAdvertiseQuery() (*advertiseQuery.GetManyAdvertiseQuery, error) {
	var postId uuid.UUID
	if req.PostId != "" {
		parsePostId, err := uuid.Parse(req.PostId)
		if err != nil {
			return nil, err
		}
		postId = parsePostId
	}

	return &advertiseQuery.GetManyAdvertiseQuery{
		PostId:       postId,
		UserEmail:    req.UserEmail,
		Status:       req.Status,
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
		FromPrice:    req.FromPrice,
		ToPrice:      req.ToPrice,
		SortBy:       req.SortBy,
		IsDescending: req.IsDescending,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}
