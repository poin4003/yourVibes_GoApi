package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToUserReportModel(userReport *user_entity.UserReport) *models.UserReport {
	ur := &models.UserReport{
		Reason:    userReport.Reason,
		Status:    userReport.Status,
		CreatedAt: userReport.CreatedAt,
		UpdatedAt: userReport.UpdatedAt,
	}
	ur.UserId = userReport.UserId
	ur.ReportedUserId = userReport.ReportedUserId

	return ur
}

func FromUserReportModel(ur *models.UserReport) *user_entity.UserReport {
	var user = &user_entity.UserForReport{
		ID:           ur.User.ID,
		FamilyName:   ur.User.FamilyName,
		Name:         ur.User.Name,
		Email:        ur.User.Email,
		Password:     ur.User.Password,
		PhoneNumber:  ur.User.PhoneNumber,
		Birthday:     ur.User.Birthday,
		AvatarUrl:    ur.User.AvatarUrl,
		CapwallUrl:   ur.User.CapwallUrl,
		Privacy:      ur.User.Privacy,
		Biography:    ur.User.Biography,
		AuthType:     ur.User.AuthType,
		AuthGoogleId: ur.User.AuthGoogleId,
		PostCount:    ur.User.PostCount,
		FriendCount:  ur.User.FriendCount,
		Status:       ur.User.Status,
		CreatedAt:    ur.User.CreatedAt,
		UpdatedAt:    ur.User.UpdatedAt,
	}

	var reportedUser = &user_entity.UserForReport{
		ID:           ur.ReportedUser.ID,
		FamilyName:   ur.ReportedUser.FamilyName,
		Name:         ur.ReportedUser.Name,
		Email:        ur.ReportedUser.Email,
		Password:     ur.ReportedUser.Password,
		PhoneNumber:  ur.ReportedUser.PhoneNumber,
		Birthday:     ur.ReportedUser.Birthday,
		AvatarUrl:    ur.ReportedUser.AvatarUrl,
		CapwallUrl:   ur.ReportedUser.CapwallUrl,
		Privacy:      ur.ReportedUser.Privacy,
		Biography:    ur.ReportedUser.Biography,
		AuthType:     ur.ReportedUser.AuthType,
		AuthGoogleId: ur.ReportedUser.AuthGoogleId,
		PostCount:    ur.ReportedUser.PostCount,
		FriendCount:  ur.ReportedUser.FriendCount,
		Status:       ur.ReportedUser.Status,
		CreatedAt:    ur.ReportedUser.CreatedAt,
		UpdatedAt:    ur.ReportedUser.UpdatedAt,
	}

	var admin *user_entity.Admin
	if ur.AdminId != nil {
		admin = &user_entity.Admin{
			ID:          ur.Admin.ID,
			FamilyName:  ur.Admin.FamilyName,
			Name:        ur.Admin.Name,
			Email:       ur.Admin.Email,
			PhoneNumber: ur.Admin.PhoneNumber,
			IdentityId:  ur.Admin.IdentityId,
			Birthday:    ur.Admin.Birthday,
			Status:      ur.Admin.Status,
			Role:        ur.Admin.Role,
			CreatedAt:   ur.Admin.CreatedAt,
			UpdatedAt:   ur.Admin.UpdatedAt,
		}
	}

	var userReport = &user_entity.UserReport{
		AdminId:      ur.AdminId,
		User:         user,
		ReportedUser: reportedUser,
		Admin:        admin,
		Reason:       ur.Reason,
		Status:       ur.Status,
		CreatedAt:    ur.CreatedAt,
		UpdatedAt:    ur.UpdatedAt,
	}

	userReport.UserId = ur.UserId
	userReport.ReportedUserId = ur.ReportedUserId

	return userReport
}
