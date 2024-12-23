package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
)

func NewCommentReportResult(
	commentReport *comment_entity.CommentReport,
) *common.CommentReportResult {
	var user = &common.UserForReportResult{
		ID:          commentReport.User.ID,
		FamilyName:  commentReport.User.FamilyName,
		Name:        commentReport.User.Name,
		Email:       commentReport.User.Email,
		PhoneNumber: commentReport.User.PhoneNumber,
		Birthday:    commentReport.User.Birthday,
		AvatarUrl:   commentReport.User.AvatarUrl,
		CapwallUrl:  commentReport.User.CapwallUrl,
		Privacy:     commentReport.User.Privacy,
		Biography:   commentReport.User.Biography,
		PostCount:   commentReport.User.PostCount,
		FriendCount: commentReport.User.FriendCount,
		Status:      commentReport.User.Status,
		CreatedAt:   commentReport.User.CreatedAt,
		UpdatedAt:   commentReport.User.UpdatedAt,
	}

	var parentPost *common.PostForReportResult
	if commentReport.Post.ParentPost != nil {
		var medias []*common.MediaResult
		for _, media := range commentReport.Post.ParentPost.Media {
			medias = append(medias, &common.MediaResult{
				ID:        media.ID,
				MediaUrl:  media.MediaUrl,
				PostId:    media.PostId,
				Status:    media.Status,
				CreatedAt: media.CreatedAt,
				UpdatedAt: media.UpdatedAt,
			})
		}
		parentPost = &common.PostForReportResult{
			ID:              commentReport.Post.ParentPost.ID,
			UserId:          commentReport.Post.ParentPost.UserId,
			User:            NewUserForReportResult(commentReport.Post.ParentPost.User),
			ParentId:        commentReport.Post.ParentPost.ParentId,
			Content:         commentReport.Post.ParentPost.Content,
			LikeCount:       commentReport.Post.ParentPost.LikeCount,
			CommentCount:    commentReport.Post.ParentPost.CommentCount,
			Privacy:         commentReport.Post.ParentPost.Privacy,
			Location:        commentReport.Post.ParentPost.Location,
			IsAdvertisement: commentReport.Post.ParentPost.IsAdvertisement,
			Status:          commentReport.Post.ParentPost.Status,
			CreatedAt:       commentReport.Post.ParentPost.CreatedAt,
			UpdatedAt:       commentReport.Post.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*common.MediaResult
	for _, media := range commentReport.Post.Media {
		medias = append(medias, &common.MediaResult{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	post := &common.PostForReportResult{
		ID:              commentReport.Post.ID,
		UserId:          commentReport.Post.UserId,
		User:            NewUserForReportResult(commentReport.Post.User),
		ParentId:        commentReport.Post.ParentId,
		ParentPost:      parentPost,
		Content:         commentReport.Post.Content,
		LikeCount:       commentReport.Post.LikeCount,
		CommentCount:    commentReport.Post.CommentCount,
		Privacy:         commentReport.Post.Privacy,
		Location:        commentReport.Post.Location,
		IsAdvertisement: commentReport.Post.IsAdvertisement,
		Status:          commentReport.Post.Status,
		CreatedAt:       commentReport.Post.CreatedAt,
		UpdatedAt:       commentReport.Post.UpdatedAt,
		Media:           medias,
	}

	reportedComment := &common.CommentForReportResult{
		ID:              commentReport.ReportedComment.ID,
		PostId:          commentReport.ReportedComment.PostId,
		UserId:          commentReport.ReportedComment.UserId,
		User:            NewUserForReportResult(commentReport.ReportedComment.User),
		ParentId:        commentReport.ReportedComment.ParentId,
		Content:         commentReport.ReportedComment.Content,
		LikeCount:       commentReport.ReportedComment.LikeCount,
		RepCommentCount: commentReport.ReportedComment.RepCommentCount,
		CreatedAt:       commentReport.ReportedComment.CreatedAt,
		UpdatedAt:       commentReport.ReportedComment.UpdatedAt,
		Status:          commentReport.ReportedComment.Status,
	}

	var admin *common.AdminResult
	if commentReport.AdminId != nil {
		admin = &common.AdminResult{
			ID:          commentReport.Admin.ID,
			FamilyName:  commentReport.Admin.FamilyName,
			Name:        commentReport.Admin.Name,
			Email:       commentReport.Admin.Email,
			PhoneNumber: commentReport.Admin.PhoneNumber,
			IdentityId:  commentReport.Admin.IdentityId,
			Birthday:    commentReport.Admin.Birthday,
			Status:      commentReport.Admin.Status,
			Role:        commentReport.Admin.Role,
			CreatedAt:   commentReport.Admin.CreatedAt,
			UpdatedAt:   commentReport.Admin.UpdatedAt,
		}
	}

	var commentReportResult = &common.CommentReportResult{
		AdminId:         commentReport.AdminId,
		User:            user,
		ReportedComment: reportedComment,
		Post:            post,
		Admin:           admin,
		Reason:          commentReport.Reason,
		Status:          commentReport.Status,
		CreatedAt:       commentReport.CreatedAt,
		UpdatedAt:       commentReport.UpdatedAt,
	}
	commentReportResult.UserId = commentReport.UserId
	commentReportResult.ReportedCommentId = commentReport.ReportedCommentId

	return commentReportResult
}

func NewCommentReportShortVerResult(
	commentReport *comment_entity.CommentReport,
) *common.CommentReportShortVerResult {
	if commentReport == nil {
		return nil
	}

	var adminEmail *string
	if commentReport.Admin != nil {
		adminEmail = &commentReport.Admin.Email
	}

	var commentReportResult = &common.CommentReportShortVerResult{
		AdminId:    commentReport.AdminId,
		Reason:     commentReport.Reason,
		UserEmail:  commentReport.User.Email,
		AdminEmail: adminEmail,
		Status:     commentReport.Status,
		CreatedAt:  commentReport.CreatedAt,
		UpdatedAt:  commentReport.UpdatedAt,
	}
	commentReportResult.UserId = commentReport.UserId
	commentReportResult.ReportedCommentId = commentReport.ReportedCommentId

	return commentReportResult
}
