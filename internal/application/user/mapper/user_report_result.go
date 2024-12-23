package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/common"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
)

func NewUserReportResult(
	userReport *user_entity.UserReport,
) *common.UserReportResult {
	if userReport == nil {
		return nil
	}

	var user = &common.UserForReportResult{
		ID:          userReport.User.ID,
		FamilyName:  userReport.User.FamilyName,
		Name:        userReport.User.Name,
		Email:       userReport.User.Email,
		PhoneNumber: userReport.User.PhoneNumber,
		Birthday:    userReport.User.Birthday,
		AvatarUrl:   userReport.User.AvatarUrl,
		CapwallUrl:  userReport.User.CapwallUrl,
		Privacy:     userReport.User.Privacy,
		Biography:   userReport.User.Biography,
		PostCount:   userReport.User.PostCount,
		FriendCount: userReport.User.FriendCount,
		Status:      userReport.User.Status,
		CreatedAt:   userReport.User.CreatedAt,
		UpdatedAt:   userReport.User.UpdatedAt,
	}

	var reportedUser = &common.UserForReportResult{
		ID:          userReport.ReportedUser.ID,
		FamilyName:  userReport.ReportedUser.FamilyName,
		Name:        userReport.ReportedUser.Name,
		Email:       userReport.ReportedUser.Email,
		PhoneNumber: userReport.ReportedUser.PhoneNumber,
		Birthday:    userReport.ReportedUser.Birthday,
		AvatarUrl:   userReport.ReportedUser.AvatarUrl,
		CapwallUrl:  userReport.ReportedUser.CapwallUrl,
		Privacy:     userReport.ReportedUser.Privacy,
		Biography:   userReport.ReportedUser.Biography,
		PostCount:   userReport.ReportedUser.PostCount,
		FriendCount: userReport.ReportedUser.FriendCount,
		Status:      userReport.ReportedUser.Status,
		CreatedAt:   userReport.ReportedUser.CreatedAt,
		UpdatedAt:   userReport.ReportedUser.UpdatedAt,
	}

	var admin *common.AdminResult
	if userReport.AdminId != nil {
		admin = &common.AdminResult{
			ID:          userReport.Admin.ID,
			FamilyName:  userReport.Admin.FamilyName,
			Name:        userReport.Admin.Name,
			Email:       userReport.Admin.Email,
			PhoneNumber: userReport.Admin.PhoneNumber,
			IdentityId:  userReport.Admin.IdentityId,
			Birthday:    userReport.Admin.Birthday,
			Status:      userReport.Admin.Status,
			Role:        userReport.Admin.Role,
			CreatedAt:   userReport.Admin.CreatedAt,
			UpdatedAt:   userReport.Admin.UpdatedAt,
		}
	}

	var userReportResult = &common.UserReportResult{
		AdminId:      userReport.AdminId,
		User:         user,
		ReportedUser: reportedUser,
		Admin:        admin,
		Reason:       userReport.Reason,
		Status:       userReport.Status,
		CreatedAt:    userReport.CreatedAt,
		UpdatedAt:    userReport.UpdatedAt,
	}
	userReportResult.UserId = userReport.UserId
	userReportResult.ReportedUserId = userReport.ReportedUserId

	return userReportResult
}

func NewUserReportShortVerResult(
	userReport *user_entity.UserReport,
) *common.UserReportShortVerResult {
	if userReport == nil {
		return nil
	}

	var adminEmail *string
	if userReport.Admin != nil {
		adminEmail = &userReport.Admin.Email
	}

	var userReportResult = &common.UserReportShortVerResult{
		AdminId:           userReport.AdminId,
		UserEmail:         userReport.User.Email,
		ReportedUserEmail: userReport.ReportedUser.Email,
		AdminEmail:        adminEmail,
		Reason:            userReport.Reason,
		Status:            userReport.Status,
		CreatedAt:         userReport.CreatedAt,
		UpdatedAt:         userReport.UpdatedAt,
	}
	userReportResult.UserId = userReport.UserId
	userReportResult.ReportedUserId = userReport.ReportedUserId

	return userReportResult
}
