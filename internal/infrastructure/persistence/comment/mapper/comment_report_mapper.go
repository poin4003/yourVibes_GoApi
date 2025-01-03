package mapper

import (
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToCommentReportModel(commentReport *commentEntity.CommentReport) *models.CommentReport {
	cr := &models.CommentReport{
		Reason:    commentReport.Reason,
		Status:    commentReport.Status,
		CreatedAt: commentReport.CreatedAt,
		UpdatedAt: commentReport.UpdatedAt,
	}
	cr.UserId = commentReport.UserId
	cr.ReportedCommentId = commentReport.ReportedCommentId

	return cr
}

func FromCommentReportModel(cr *models.CommentReport) *commentEntity.CommentReport {
	var user = &commentEntity.UserForReport{
		ID:           cr.User.ID,
		FamilyName:   cr.User.FamilyName,
		Name:         cr.User.Name,
		Email:        cr.User.Email,
		Password:     cr.User.Password,
		PhoneNumber:  cr.User.PhoneNumber,
		Birthday:     cr.User.Birthday,
		AvatarUrl:    cr.User.AvatarUrl,
		CapwallUrl:   cr.User.CapwallUrl,
		Privacy:      cr.User.Privacy,
		Biography:    cr.User.Biography,
		AuthType:     cr.User.AuthType,
		AuthGoogleId: cr.User.AuthGoogleId,
		PostCount:    cr.User.PostCount,
		FriendCount:  cr.User.FriendCount,
		Status:       cr.User.Status,
		CreatedAt:    cr.User.CreatedAt,
		UpdatedAt:    cr.User.UpdatedAt,
	}

	var parentPost *commentEntity.PostForReport
	if cr.ReportedComment.Post.ParentPost != nil {
		var medias []*commentEntity.Media
		for _, media := range cr.ReportedComment.Post.ParentPost.Media {
			medias = append(medias, &commentEntity.Media{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &commentEntity.PostForReport{
			ID:              cr.ReportedComment.Post.ParentPost.ID,
			UserId:          cr.ReportedComment.Post.ParentPost.UserId,
			User:            ToUserForReportEntity(&cr.ReportedComment.Post.ParentPost.User),
			ParentId:        cr.ReportedComment.Post.ParentPost.ParentId,
			Content:         cr.ReportedComment.Post.ParentPost.Content,
			LikeCount:       cr.ReportedComment.Post.ParentPost.LikeCount,
			CommentCount:    cr.ReportedComment.Post.ParentPost.CommentCount,
			Privacy:         cr.ReportedComment.Post.ParentPost.Privacy,
			Location:        cr.ReportedComment.Post.ParentPost.Location,
			IsAdvertisement: cr.ReportedComment.Post.ParentPost.IsAdvertisement,
			Status:          cr.ReportedComment.Post.ParentPost.Status,
			CreatedAt:       cr.ReportedComment.Post.ParentPost.CreatedAt,
			UpdatedAt:       cr.ReportedComment.Post.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*commentEntity.Media
	for _, media := range cr.ReportedComment.Post.Media {
		medias = append(medias, &commentEntity.Media{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	post := &commentEntity.PostForReport{
		ID:              cr.ReportedComment.Post.ID,
		UserId:          cr.ReportedComment.Post.UserId,
		User:            ToUserForReportEntity(&cr.ReportedComment.Post.User),
		ParentId:        cr.ReportedComment.Post.ParentId,
		ParentPost:      parentPost,
		Content:         cr.ReportedComment.Post.Content,
		LikeCount:       cr.ReportedComment.Post.LikeCount,
		CommentCount:    cr.ReportedComment.Post.CommentCount,
		Privacy:         cr.ReportedComment.Post.Privacy,
		Location:        cr.ReportedComment.Post.Location,
		IsAdvertisement: cr.ReportedComment.Post.IsAdvertisement,
		Status:          cr.ReportedComment.Post.Status,
		CreatedAt:       cr.ReportedComment.Post.CreatedAt,
		UpdatedAt:       cr.ReportedComment.Post.UpdatedAt,
		Media:           medias,
	}

	reportedComment := &commentEntity.CommentForReport{
		ID:              cr.ReportedComment.ID,
		PostId:          cr.ReportedComment.PostId,
		UserId:          cr.ReportedComment.UserId,
		User:            ToUserForReportEntity(&cr.ReportedComment.User),
		ParentId:        cr.ReportedComment.ParentId,
		Content:         cr.ReportedComment.Content,
		LikeCount:       cr.ReportedComment.LikeCount,
		RepCommentCount: cr.ReportedComment.RepCommentCount,
		CreatedAt:       cr.ReportedComment.CreatedAt,
		UpdatedAt:       cr.ReportedComment.UpdatedAt,
		Status:          cr.ReportedComment.Status,
	}

	var admin *commentEntity.Admin
	if cr.AdminId != nil {
		admin = &commentEntity.Admin{
			ID:          cr.Admin.ID,
			FamilyName:  cr.Admin.FamilyName,
			Name:        cr.Admin.Name,
			Email:       cr.Admin.Email,
			PhoneNumber: cr.Admin.PhoneNumber,
			IdentityId:  cr.Admin.IdentityId,
			Birthday:    cr.Admin.Birthday,
			Status:      cr.Admin.Status,
			Role:        cr.Admin.Role,
			CreatedAt:   cr.Admin.CreatedAt,
			UpdatedAt:   cr.Admin.UpdatedAt,
		}
	}

	var commentReport = &commentEntity.CommentReport{
		AdminId:         cr.AdminId,
		User:            user,
		ReportedComment: reportedComment,
		Post:            post,
		Admin:           admin,
		Reason:          cr.Reason,
		Status:          cr.Status,
		CreatedAt:       cr.CreatedAt,
		UpdatedAt:       cr.UpdatedAt,
	}
	commentReport.UserId = cr.UserId
	commentReport.ReportedCommentId = cr.ReportedCommentId

	return commentReport
}

func ToUserForReportEntity(
	userModel *models.User,
) *commentEntity.UserForReport {
	if userModel == nil {
		return nil
	}

	var userForReport = &commentEntity.UserForReport{
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
