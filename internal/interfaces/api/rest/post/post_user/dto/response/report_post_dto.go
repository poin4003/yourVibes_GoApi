package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"time"
)

type PostReportDto struct {
	UserId         uuid.UUID         `json:"user_id"`
	ReportedPostId uuid.UUID         `json:"reported_post_id"`
	User           *UserForReportDto `json:"user"`
	ReportedPost   *PostForReportDto `json:"reported_post"`
	Reason         string            `json:"reason"`
	Status         bool              `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func ToPostReportDto(postResult *common.PostReportResult) *PostReportDto {
	return &PostReportDto{
		UserId:         postResult.UserId,
		ReportedPostId: postResult.ReportedPostId,
		User:           ToUserForReportDto(postResult.User),
		ReportedPost:   ToPostForReportDto(postResult.ReportedPost),
		Reason:         postResult.Reason,
		Status:         postResult.Status,
		CreatedAt:      postResult.CreatedAt,
		UpdatedAt:      postResult.UpdatedAt,
	}
}
