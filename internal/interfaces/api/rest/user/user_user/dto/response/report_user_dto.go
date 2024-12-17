package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	"time"
)

type UserReportDto struct {
	UserId         uuid.UUID         `json:"user_id"`
	ReportedUserId uuid.UUID         `json:"reported_user_id"`
	User           *UserForReportDto `json:"user"`
	ReportedUser   *UserForReportDto `json:"reported_user"`
	Reason         string            `json:"reason"`
	Status         bool              `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func ToUserReportDto(userResult *common.UserReportResult) *UserReportDto {
	return &UserReportDto{
		UserId:         userResult.UserId,
		ReportedUserId: userResult.ReportedUserId,
		User:           ToUserForReportDto(userResult.User),
		ReportedUser:   ToUserForReportDto(userResult.ReportedUser),
		Reason:         userResult.Reason,
		Status:         userResult.Status,
		CreatedAt:      userResult.CreatedAt,
		UpdatedAt:      userResult.UpdatedAt,
	}
}
