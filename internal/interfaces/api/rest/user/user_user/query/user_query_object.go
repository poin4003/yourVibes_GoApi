package query

import (
	"github.com/google/uuid"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"time"
)

type UserQueryObject struct {
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

func (req *UserQueryObject) ToGetOneUserQuery(
	userId uuid.UUID,
	authenticatedUserId uuid.UUID,
) (*user_query.GetOneUserQuery, error) {
	return &user_query.GetOneUserQuery{
		UserId:              userId,
		AuthenticatedUserId: authenticatedUserId,
	}, nil
}

func (req *UserQueryObject) ToGetManyUserQuery() (*user_query.GetManyUserQuery, error) {
	return &user_query.GetManyUserQuery{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Birthday:     req.Birthday,
		CreatedAt:    req.CreatedAt,
		SortBy:       req.SortBy,
		IsDescending: req.IsDescending,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}
