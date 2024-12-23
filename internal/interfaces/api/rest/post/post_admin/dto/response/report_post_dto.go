package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
)

type PostReportDto struct {
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

func ToPostReportDto(postResult *common.PostReportResult) *PostReportDto {
	if postResult == nil {
		return nil
	}

	return &PostReportDto{
		UserId:         postResult.UserId,
		ReportedPostId: postResult.ReportedPostId,
		AdminId:        postResult.AdminId,
		User:           ToUserForReportDto(postResult.User),
		ReportedPost:   ToPostForReportDto(postResult.ReportedPost),
		Admin:          ToAdminDto(postResult.Admin),
		Reason:         postResult.Reason,
		Status:         postResult.Status,
		CreatedAt:      postResult.CreatedAt,
		UpdatedAt:      postResult.UpdatedAt,
	}
}

func ToPostReportShortVerDto(
	postResult *common.PostReportShortVerResult,
) *PostReportShortVerDto {
	if postResult == nil {
		return nil
	}

	return &PostReportShortVerDto{
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
