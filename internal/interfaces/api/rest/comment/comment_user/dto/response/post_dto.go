package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type PostForReportDto struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *UserForReportDto
	ParentId        *uuid.UUID
	ParentPost      *PostForReportDto
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*MediaDto
}

func ToPostForReportDto(postResult *common.PostForReportResult) *PostForReportDto {
	var parentPost *PostForReportDto

	if postResult.ParentPost != nil {
		parentPost = &PostForReportDto{
			ID:              postResult.ParentPost.ID,
			UserId:          postResult.ParentPost.UserId,
			User:            ToUserForReportDto(postResult.ParentPost.User),
			Content:         postResult.ParentPost.Content,
			LikeCount:       postResult.ParentPost.LikeCount,
			CommentCount:    postResult.ParentPost.CommentCount,
			Privacy:         postResult.ParentPost.Privacy,
			Location:        postResult.ParentPost.Location,
			IsAdvertisement: postResult.ParentPost.IsAdvertisement,
			Status:          postResult.ParentPost.Status,
			CreatedAt:       postResult.ParentPost.CreatedAt,
			UpdatedAt:       postResult.ParentPost.UpdatedAt,
			Media:           ToMediaDto(postResult.ParentPost.Media),
		}
	}

	return &PostForReportDto{
		ID:              postResult.ID,
		UserId:          postResult.UserId,
		User:            ToUserForReportDto(postResult.User),
		ParentId:        postResult.ParentId,
		ParentPost:      parentPost,
		Content:         postResult.Content,
		LikeCount:       postResult.LikeCount,
		CommentCount:    postResult.CommentCount,
		Privacy:         postResult.Privacy,
		Location:        postResult.Location,
		IsAdvertisement: postResult.IsAdvertisement,
		Status:          postResult.Status,
		CreatedAt:       postResult.CreatedAt,
		UpdatedAt:       postResult.UpdatedAt,
		Media:           ToMediaDto(postResult.Media),
	}
}
