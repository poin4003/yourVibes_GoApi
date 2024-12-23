package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	post_query "github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
)

type PostReportQueryObject struct {
	Reason       string    `form:"reason,omitempty"`
	UserEmail    string    `form:"user_email,omitempty"`
	AdminEmail   string    `form:"admin_email,omitemtpy"`
	FromDate     time.Time `form:"from_date,omitemtpy"`
	ToDate       time.Time `form:"to_date,omitemtpy"`
	CreatedAt    time.Time `form:"created_at,omitempty"`
	Status       *bool     `form:"status,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"is_descending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}

func ValidatePostReportQueryObject(input interface{}) error {
	query, ok := input.(*PostReportQueryObject)
	if !ok {
		return fmt.Errorf("validatePostReportQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.Reason, validation.Length(10, 255)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func ToGetOnePostReportQuery(
	userId uuid.UUID,
	reportedPostId uuid.UUID,
) (*post_query.GetOnePostReportQuery, error) {
	return &post_query.GetOnePostReportQuery{
		UserId:         userId,
		ReportedPostId: reportedPostId,
	}, nil
}

func (req *PostReportQueryObject) ToGetManyPostQuery() (*post_query.GetManyPostReportQuery, error) {
	return &post_query.GetManyPostReportQuery{
		Reason:       req.Reason,
		UserEmail:    req.UserEmail,
		AdminEmail:   req.AdminEmail,
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
		CreatedAt:    req.CreatedAt,
		Status:       req.Status,
		SortBy:       req.SortBy,
		IsDescending: req.IsDescending,
		Limit:        req.Limit,
		Page:         req.Page,
	}, nil
}
