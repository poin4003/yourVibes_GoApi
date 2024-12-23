package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
)

type CommentReportQueryObject struct {
	Reason       string    `form:"reason,omitempty"`
	UserEmail    string    `form:"user_email,omitempty"`
	AdminEmail   string    `form:"admin_email,omitempty"`
	FromDate     time.Time `form:"from_date,omitempty"`
	ToDate       time.Time `form:"to_date,omitempty"`
	CreatedAt    time.Time `form:"created_at,omitempty"`
	Status       *bool     `form:"status,omitempty"`
	SortBy       string    `form:"sort_by,omitempty"`
	IsDescending bool      `form:"is_descending,omitempty"`
	Limit        int       `form:"limit,omitempty"`
	Page         int       `form:"page,omitempty"`
}

func ValidateCommentReportQueryObject(input interface{}) error {
	query, ok := input.(*CommentReportQueryObject)
	if !ok {
		return fmt.Errorf("validateCommentReportQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.UserEmail, is.Email),
		validation.Field(&query.AdminEmail, is.Email),
		validation.Field(&query.Reason, validation.Length(10, 255)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func ToGetOneCommentReportQuery(
	userId uuid.UUID,
	reportedCommentId uuid.UUID,
) (*comment_query.GetOneCommentReportQuery, error) {
	return &comment_query.GetOneCommentReportQuery{
		UserId:            userId,
		ReportedCommentId: reportedCommentId,
	}, nil
}

func (req *CommentReportQueryObject) ToGetManyCommentQuery() (*comment_query.GetManyCommentReportQuery, error) {
	return &comment_query.GetManyCommentReportQuery{
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
