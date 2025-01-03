package query

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	adminQuery "github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	"regexp"
	"time"
)

type AdminQueryObject struct {
	Name         string    `form:"name,omitempty"`
	Email        string    `form:"email,omitempty"`
	PhoneNumber  string    `form:"phone_number,omitempty"`
	IdentityId   string    `form:"identity_id,omitempty"`
	Birthday     time.Time `form:"birthday,omitempty"`
	CreatedAt    time.Time `form:"created_at,omitempty"`
	Status       *bool     `form:"status,omitempty"`
	Role         *bool     `form:"role,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"is_descending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}

func ValidateAdminQueryObject(input interface{}) error {
	query, ok := input.(*AdminQueryObject)
	if !ok {
		return fmt.Errorf("validateAdminQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Name, validation.Length(1, 510)),
		validation.Field(&query.Email, is.Email),
		validation.Field(&query.PhoneNumber, validation.Length(10, 14), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&query.IdentityId, validation.Length(10, 15), validation.Match((regexp.MustCompile((`^\d+$`))))),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *AdminQueryObject) ToGetOneAdminQuery(
	adminID uuid.UUID,
) (*adminQuery.GetOneAdminQuery, error) {
	return &adminQuery.GetOneAdminQuery{
		AdminId: adminID,
	}, nil
}

func (req *AdminQueryObject) ToGetManyAdminQuery() (*adminQuery.GetManyAdminQuery, error) {
	return &adminQuery.GetManyAdminQuery{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		IdentityId:   req.IdentityId,
		Birthday:     req.Birthday,
		CreatedAt:    req.CreatedAt,
		SortBy:       req.SortBy,
		Status:       req.Status,
		Role:         req.Role,
		IsDescending: req.IsDescending,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}
