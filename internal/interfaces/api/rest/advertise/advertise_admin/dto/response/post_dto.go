package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type PostForAdvertiseDto struct {
	ID              uuid.UUID            `json:"id"`
	UserId          uuid.UUID            `json:"user_id"`
	User            *UserForAdvertiseDto `json:"user"`
	ParentId        *uuid.UUID           `json:"parent_id"`
	ParentPost      *PostForAdvertiseDto `json:"parent_post"`
	Content         string               `json:"content"`
	LikeCount       int                  `json:"like_count"`
	CommentCount    int                  `json:"comment_count"`
	Privacy         consts.PrivacyLevel  `json:"privacy"`
	Location        string               `json:"location"`
	IsAdvertisement bool                 `json:"is_advertisement"`
	Status          bool                 `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Media           []*MediaDto          `json:"media"`
}

func ToPostForAdvertiseDto(postResult *common.PostForAdvertiseResult) *PostForAdvertiseDto {
	var parentPost *PostForAdvertiseDto

	if postResult.ParentPost != nil {
		parentPost = &PostForAdvertiseDto{
			ID:              postResult.ParentPost.ID,
			UserId:          postResult.ParentPost.UserId,
			User:            ToUserForAdvertiseDto(postResult.ParentPost.User),
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

	return &PostForAdvertiseDto{
		ID:              postResult.ID,
		UserId:          postResult.UserId,
		User:            ToUserForAdvertiseDto(postResult.User),
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
