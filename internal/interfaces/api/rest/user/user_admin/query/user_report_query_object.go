package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	user_query "github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
)

type UserReportQueryObject struct {
	Reason            string    `form:"reason,omitempty"`
	UserEmail         string    `form:"user_email,omitempty"`
	ReportedUserEmail string    `form:"reported_user_email,omitempty"`
	AdminEmail        string    `form:"admin_email,omitemtpy"`
	FromDate          time.Time `form:"from_date"`
	ToDate            time.Time `form:"to_date"`
	CreatedAt         time.Time `form:"created_at,omitempty"`
	Status            *bool     `form:"status,omitempty"`
	SortBy            string    `form:"sort_by,omitempty"`
	IsDescending      bool      `form:"is_descending,omitempty"`
	Limit             int       `form:"limit,omitempty"`
	Page              int       `form:"page,omitempty"`
}

func ValidateUserReportQueryObject(input interface{}) error {
	query, ok := input.(*UserReportQueryObject)
	if !ok {
		return fmt.Errorf("validateUserReportQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Reason, validation.Length(10, 255)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func ToGetOneUserReportQuery(
	userId uuid.UUID,
	reportedUserId uuid.UUID,
) (*user_query.GetOneUserReportQuery, error) {
	return &user_query.GetOneUserReportQuery{
		UserId:         userId,
		ReportedUserId: reportedUserId,
	}, nil
}

func (req *UserReportQueryObject) ToGetManyUserQuery() (*user_query.GetManyUserReportQuery, error) {
	return &user_query.GetManyUserReportQuery{
		Reason:            req.Reason,
		UserEmail:         req.UserEmail,
		ReportedUserEmail: req.ReportedUserEmail,
		AdminEmail:        req.AdminEmail,
		FromDate:          req.FromDate,
		ToDate:            req.ToDate,
		CreatedAt:         req.CreatedAt,
		Status:            req.Status,
		SortBy:            req.SortBy,
		IsDescending:      req.IsDescending,
		Limit:             req.Limit,
		Page:              req.Page,
	}, nil
}
