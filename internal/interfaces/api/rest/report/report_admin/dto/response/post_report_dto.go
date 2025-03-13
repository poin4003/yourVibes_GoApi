package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/report/common"
)

type PostReportDto struct {
	ReportId       uuid.UUID         `json:"report_id"`
	UserId         uuid.UUID         `json:"user_id"`
	ReportedPostId uuid.UUID         `json:"reported_post_id"`
	AdminId        *uuid.UUID        `json:"admin_id"`
	User           *UserForReportDto `json:"user"`
	ReportedPost   *PostForReportDto `json:"reported_post"`
	Admin          *AdminDto         `json:"admin"`
	Reason         string            `json:"reason"`
	Status         bool              `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type PostReportShortVerDto struct {
	ReportId       uuid.UUID  `json:"report_id"`
	UserId         uuid.UUID  `json:"user_id"`
	ReportedPostId uuid.UUID  `json:"reported_post_id"`
	UserEmail      string     `json:"user_email"`
	AdminEmail     *string    `json:"admin_email"`
	AdminId        *uuid.UUID `json:"admin_id"`
	Reason         string     `json:"reason"`
	Status         bool       `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func ToPostReportDto(postReportResult *common.PostReportResult) *PostReportDto {
	if postReportResult == nil {
		return nil
	}

	return &PostReportDto{
		ReportId:       postReportResult.ReportId,
		UserId:         postReportResult.UserId,
		ReportedPostId: postReportResult.ReportedPostId,
		AdminId:        postReportResult.AdminId,
		User:           ToUserForReportDto(postReportResult.User),
		ReportedPost:   ToPostForReportDto(postReportResult.ReportedPost),
		Admin:          ToAdminDto(postReportResult.Admin),
		Reason:         postReportResult.Reason,
		Status:         postReportResult.Status,
		CreatedAt:      postReportResult.CreatedAt,
		UpdatedAt:      postReportResult.UpdatedAt,
	}
}

func ToPostReportShortVerDto(
	postResult *common.PostReportShortVerResult,
) *PostReportShortVerDto {
	if postResult == nil {
		return nil
	}

	return &PostReportShortVerDto{
		ReportId:       postResult.ReportId,
		UserId:         postResult.UserId,
		ReportedPostId: postResult.ReportedPostId,
		UserEmail:      postResult.UserEmail,
		AdminEmail:     postResult.AdminEmail,
		AdminId:        postResult.AdminId,
		Reason:         postResult.Reason,
		Status:         postResult.Status,
		CreatedAt:      postResult.CreatedAt,
		UpdatedAt:      postResult.UpdatedAt,
	}
}
