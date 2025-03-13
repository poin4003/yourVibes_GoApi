package query

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	reportQuery "github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type ReportDetailQueryObject struct {
	ReportType consts.ReportType `form:"report_type"`
}

type ReportQueryObject struct {
	ReportType   consts.ReportType `form:"report_type"`
	Reason       string            `form:"reason,omitempty"`
	UserEmail    string            `form:"user_email,omitempty"`
	AdminEmail   string            `form:"admin_email,omitempty"`
	FromDate     time.Time         `form:"from_date,omitempty"`
	ToDate       time.Time         `form:"to_date,omitempty"`
	CreatedAt    time.Time         `form:"created_at,omitempty"`
	Status       *bool             `form:"status,omitempty"`
	SortBy       string            `form:"sort_by,omitempty"`
	IsDescending bool              `form:"is_descending,omitempty"`
	Limit        int               `form:"limit,omitempty"`
	Page         int               `form:"page,omitempty"`
}

func ValidateReportDetailQueryObject(input interface{}) error {
	query, ok := input.(*ReportDetailQueryObject)
	if !ok {
		return fmt.Errorf("validateReportDetailQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.ReportType, validation.In(consts.ReportTypes...)),
	)
}

func ValidateReportQueryObject(input interface{}) error {
	query, ok := input.(*ReportQueryObject)
	if !ok {
		return fmt.Errorf("validateReportQueryObject failed")
	}

	return validation.ValidateStruct(query,
		validation.Field(&query.ReportType, validation.In(consts.ReportTypes...)),
		validation.Field(&query.UserEmail, is.Email),
		validation.Field(&query.AdminEmail, is.Email),
		validation.Field(&query.Reason, validation.Length(10, 255)),
		validation.Field(&query.Limit, validation.Min(0)),
		validation.Field(&query.Page, validation.Min(0)),
	)
}

func (req *ReportDetailQueryObject) ToGetOneReportQuery(
	reportedId uuid.UUID,
) (*reportQuery.GetOneReportQuery, error) {
	return &reportQuery.GetOneReportQuery{
		ReportType: req.ReportType,
		ReportedId: reportedId,
	}, nil
}

func (req *ReportQueryObject) ToGetManyReportQuery() (*reportQuery.GetManyReportQuery, error) {
	return &reportQuery.GetManyReportQuery{
		ReportType:   req.ReportType,
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
