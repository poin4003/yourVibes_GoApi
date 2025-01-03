package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	userQuery "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"regexp"
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

func ValidateUserQueryObject(input interface{}) error {
	query, ok := input.(*UserQueryObject)
	if !ok {
		return fmt.Errorf("validateUserQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Name, validation.Length(1, 510)),
		validation.Field(&query.Email, is.Email),
		validation.Field(&query.PhoneNumber, validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *UserQueryObject) ToGetOneUserQuery(
	userId uuid.UUID,
	authenticatedUserId uuid.UUID,
) (*userQuery.GetOneUserQuery, error) {
	return &userQuery.GetOneUserQuery{
		UserId:              userId,
		AuthenticatedUserId: authenticatedUserId,
	}, nil
}

func (req *UserQueryObject) ToGetManyUserQuery() (*userQuery.GetManyUserQuery, error) {
	return &userQuery.GetManyUserQuery{
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
