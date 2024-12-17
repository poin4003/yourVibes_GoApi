package mapper

import (
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToPostReportModel(postReport *post_entity.PostReport) *models.PostReport {
	pr := &models.PostReport{
		Reason:    postReport.Reason,
		Status:    postReport.Status,
		CreatedAt: postReport.CreatedAt,
		UpdatedAt: postReport.UpdatedAt,
	}
	pr.UserId = postReport.UserId
	pr.ReportedPostId = postReport.ReportedPostId

	return pr
}

func FromPostReportModel(pr *models.PostReport) *post_entity.PostReport {
	var user = &post_entity.UserForReport{
		ID:           pr.User.ID,
		FamilyName:   pr.User.FamilyName,
		Name:         pr.User.Name,
		Email:        pr.User.Email,
		Password:     pr.User.Password,
		PhoneNumber:  pr.User.PhoneNumber,
		Birthday:     pr.User.Birthday,
		AvatarUrl:    pr.User.AvatarUrl,
		CapwallUrl:   pr.User.CapwallUrl,
		Privacy:      pr.User.Privacy,
		Biography:    pr.User.Biography,
		AuthType:     pr.User.AuthType,
		AuthGoogleId: pr.User.AuthGoogleId,
		PostCount:    pr.User.PostCount,
		FriendCount:  pr.User.FriendCount,
		Status:       pr.User.Status,
		CreatedAt:    pr.User.CreatedAt,
		UpdatedAt:    pr.User.UpdatedAt,
	}

	var parentPost *post_entity.PostForReport
	if pr.ReportedPost.ParentPost != nil {
		var medias []*post_entity.Media
		for _, media := range pr.ReportedPost.ParentPost.Media {
			medias = append(medias, &post_entity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &post_entity.PostForReport{
			ID:              pr.ReportedPost.ParentPost.ID,
			UserId:          pr.ReportedPost.ParentPost.UserId,
			User:            ToUserForReportEntity(&pr.ReportedPost.ParentPost.User),
			ParentId:        pr.ReportedPost.ParentPost.ParentId,
			Content:         pr.ReportedPost.ParentPost.Content,
			LikeCount:       pr.ReportedPost.ParentPost.LikeCount,
			CommentCount:    pr.ReportedPost.ParentPost.CommentCount,
			Privacy:         pr.ReportedPost.ParentPost.Privacy,
			Location:        pr.ReportedPost.ParentPost.Location,
			IsAdvertisement: pr.ReportedPost.ParentPost.IsAdvertisement,
			Status:          pr.ReportedPost.ParentPost.Status,
			CreatedAt:       pr.ReportedPost.ParentPost.CreatedAt,
			UpdatedAt:       pr.ReportedPost.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*post_entity.Media
	for _, media := range pr.ReportedPost.Media {
		medias = append(medias, &post_entity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	reportedPost := &post_entity.PostForReport{
		ID:              pr.ReportedPost.ID,
		UserId:          pr.ReportedPost.UserId,
		User:            ToUserForReportEntity(&pr.ReportedPost.User),
		ParentId:        pr.ReportedPost.ParentId,
		ParentPost:      parentPost,
		Content:         pr.ReportedPost.Content,
		LikeCount:       pr.ReportedPost.LikeCount,
		CommentCount:    pr.ReportedPost.CommentCount,
		Privacy:         pr.ReportedPost.Privacy,
		Location:        pr.ReportedPost.Location,
		IsAdvertisement: pr.ReportedPost.IsAdvertisement,
		Status:          pr.ReportedPost.Status,
		CreatedAt:       pr.ReportedPost.CreatedAt,
		UpdatedAt:       pr.ReportedPost.UpdatedAt,
		Media:           medias,
	}

	var admin *post_entity.Admin
	if pr.AdminId != nil {
		admin = &post_entity.Admin{
			ID:          pr.Admin.ID,
			FamilyName:  pr.Admin.FamilyName,
			Name:        pr.Admin.Name,
			Email:       pr.Admin.Email,
			PhoneNumber: pr.Admin.PhoneNumber,
			IdentityId:  pr.Admin.IdentityId,
			Birthday:    pr.Admin.Birthday,
			Status:      pr.Admin.Status,
			Role:        pr.Admin.Role,
			CreatedAt:   pr.Admin.CreatedAt,
			UpdatedAt:   pr.Admin.UpdatedAt,
		}
	}

	var postReport = &post_entity.PostReport{
		AdminId:      pr.AdminId,
		User:         user,
		ReportedPost: reportedPost,
		Admin:        admin,
		Reason:       pr.Reason,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		UpdatedAt:    pr.UpdatedAt,
	}
	postReport.UserId = pr.UserId
	postReport.ReportedPostId = pr.ReportedPostId

	return postReport
}

func ToUserForReportEntity(
	userModel *models.User,
) *post_entity.UserForReport {
	if userModel == nil {
		return nil
	}

	var userForReport = &post_entity.UserForReport{
		FamilyName:   userModel.FamilyName,
		Name:         userModel.Name,
		Email:        userModel.Email,
		Password:     userModel.Password,
		PhoneNumber:  userModel.PhoneNumber,
		Birthday:     userModel.Birthday,
		AvatarUrl:    userModel.AvatarUrl,
		CapwallUrl:   userModel.CapwallUrl,
		Privacy:      userModel.Privacy,
		Biography:    userModel.Biography,
		AuthType:     userModel.AuthType,
		AuthGoogleId: userModel.AuthGoogleId,
		PostCount:    userModel.PostCount,
		FriendCount:  userModel.FriendCount,
		Status:       userModel.Status,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
	userForReport.ID = userModel.ID

	return userForReport
}
