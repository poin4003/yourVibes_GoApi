package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

func NewPostReportResult(
	postReport *post_entity.PostReport,
) *common.PostReportResult {
	var user = &common.UserForReportResult{
		ID:          postReport.User.ID,
		FamilyName:  postReport.User.FamilyName,
		Name:        postReport.User.Name,
		Email:       postReport.User.Email,
		PhoneNumber: postReport.User.PhoneNumber,
		Birthday:    postReport.User.Birthday,
		AvatarUrl:   postReport.User.AvatarUrl,
		CapwallUrl:  postReport.User.CapwallUrl,
		Privacy:     postReport.User.Privacy,
		Biography:   postReport.User.Biography,
		PostCount:   postReport.User.PostCount,
		FriendCount: postReport.User.FriendCount,
		Status:      postReport.User.Status,
		CreatedAt:   postReport.User.CreatedAt,
		UpdatedAt:   postReport.User.UpdatedAt,
	}

	var parentPost *common.PostForReportResult
	if postReport.ReportedPost.ParentPost != nil {
		var medias []*common.MediaResult
		for _, media := range postReport.ReportedPost.ParentPost.Media {
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
			ID:              postReport.ReportedPost.ParentPost.ID,
			UserId:          postReport.ReportedPost.ParentPost.UserId,
			User:            NewUserForReportResult(postReport.ReportedPost.ParentPost.User),
			ParentId:        postReport.ReportedPost.ParentPost.ParentId,
			Content:         postReport.ReportedPost.ParentPost.Content,
			LikeCount:       postReport.ReportedPost.ParentPost.LikeCount,
			CommentCount:    postReport.ReportedPost.ParentPost.CommentCount,
			Privacy:         postReport.ReportedPost.ParentPost.Privacy,
			Location:        postReport.ReportedPost.ParentPost.Location,
			IsAdvertisement: postReport.ReportedPost.ParentPost.IsAdvertisement,
			Status:          postReport.ReportedPost.ParentPost.Status,
			CreatedAt:       postReport.ReportedPost.ParentPost.CreatedAt,
			UpdatedAt:       postReport.ReportedPost.ParentPost.UpdatedAt,
			Media:           medias,
		}
	}

	var medias []*common.MediaResult
	for _, media := range postReport.ReportedPost.Media {
		medias = append(medias, &common.MediaResult{
			ID:        media.ID,
			MediaUrl:  media.MediaUrl,
			PostId:    media.PostId,
			Status:    media.Status,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
		})
	}

	reportedPost := &common.PostForReportResult{
		ID:              postReport.ReportedPost.ID,
		UserId:          postReport.ReportedPost.UserId,
		User:            NewUserForReportResult(postReport.ReportedPost.User),
		ParentId:        postReport.ReportedPost.ParentId,
		ParentPost:      parentPost,
		Content:         postReport.ReportedPost.Content,
		LikeCount:       postReport.ReportedPost.LikeCount,
		CommentCount:    postReport.ReportedPost.CommentCount,
		Privacy:         postReport.ReportedPost.Privacy,
		Location:        postReport.ReportedPost.Location,
		IsAdvertisement: postReport.ReportedPost.IsAdvertisement,
		Status:          postReport.ReportedPost.Status,
		CreatedAt:       postReport.ReportedPost.CreatedAt,
		UpdatedAt:       postReport.ReportedPost.UpdatedAt,
		Media:           medias,
	}

	var admin *common.AdminResult
	if postReport.AdminId != nil {
		admin = &common.AdminResult{
			ID:          postReport.Admin.ID,
			FamilyName:  postReport.Admin.FamilyName,
			Name:        postReport.Admin.Name,
			Email:       postReport.Admin.Email,
			PhoneNumber: postReport.Admin.PhoneNumber,
			IdentityId:  postReport.Admin.IdentityId,
			Birthday:    postReport.Admin.Birthday,
			Status:      postReport.Admin.Status,
			Role:        postReport.Admin.Role,
			CreatedAt:   postReport.Admin.CreatedAt,
			UpdatedAt:   postReport.Admin.UpdatedAt,
		}
	}

	var postReportResult = &common.PostReportResult{
		AdminId:      postReport.AdminId,
		User:         user,
		ReportedPost: reportedPost,
		Admin:        admin,
		Reason:       postReport.Reason,
		Status:       postReport.Status,
		CreatedAt:    postReport.CreatedAt,
		UpdatedAt:    postReport.UpdatedAt,
	}
	postReportResult.UserId = postReport.UserId
	postReportResult.ReportedPostId = postReport.ReportedPostId

	return postReportResult
}
