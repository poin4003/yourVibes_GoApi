package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
)

type UserReportDto struct {
	UserId         uuid.UUID         `json:"user_id"`
	ReportedUserId uuid.UUID         `json:"reported_user_id"`
	AdminId        *uuid.UUID        `json:"admin_id"`
	User           *UserForReportDto `json:"user"`
	ReportedUser   *UserForReportDto `json:"reported_user"`
	Admin          *AdminDto         `json:"admin"`
	Reason         string            `json:"reason"`
	Status         bool              `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type UserReportShortVerDto struct {
	UserId            uuid.UUID  `json:"user_id"`
	ReportedUserId    uuid.UUID  `json:"reported_user_id"`
	UserEmail         string     `json:"user_email"`
	ReportedUserEmail string     `json:"reported_user_email"`
	AdminEmail        *string    `json:"admin_email"`
	AdminId           *uuid.UUID `json:"admin_id"`
	Reason            string     `json:"reason"`
	Status            bool       `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func ToUserReportDto(
	userReportResult *common.UserReportResult,
) *UserReportDto {
	if userReportResult == nil {
		return nil
	}

	return &UserReportDto{
		UserId:         userReportResult.UserId,
		ReportedUserId: userReportResult.ReportedUserId,
		AdminId:        userReportResult.AdminId,
		User:           ToUserForReportDto(userReportResult.User),
		ReportedUser:   ToUserForReportDto(userReportResult.ReportedUser),
		Admin:          ToAdminDto(userReportResult.Admin),
		Reason:         userReportResult.Reason,
		Status:         userReportResult.Status,
		CreatedAt:      userReportResult.CreatedAt,
		UpdatedAt:      userReportResult.UpdatedAt,
	}
}

func ToUserReportShortVerDto(
	userReportResult *common.UserReportShortVerResult,
) *UserReportShortVerDto {
	if userReportResult == nil {
		return nil
	}

	return &UserReportShortVerDto{
		UserId:            userReportResult.UserId,
		ReportedUserId:    userReportResult.ReportedUserId,
		AdminId:           userReportResult.AdminId,
		Reason:            userReportResult.Reason,
		UserEmail:         userReportResult.UserEmail,
		ReportedUserEmail: userReportResult.ReportedUserEmail,
		AdminEmail:        userReportResult.AdminEmail,
		Status:            userReportResult.Status,
		CreatedAt:         userReportResult.CreatedAt,
		UpdatedAt:         userReportResult.UpdatedAt,
	}
}
