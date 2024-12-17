package response

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type PostDto struct {
	ID              uuid.UUID           `json:"id"`
	UserId          uuid.UUID           `json:"user_id"`
	User            *UserDto            `json:"user"`
	ParentId        *uuid.UUID          `json:"parent_id"`
	ParentPost      *PostDto            `json:"parent_post"`
	Content         string              `json:"content"`
	LikeCount       int                 `json:"like_count"`
	CommentCount    int                 `json:"comment_count"`
	Privacy         consts.PrivacyLevel `json:"privacy"`
	Location        string              `json:"location"`
	IsAdvertisement bool                `json:"is_advertisement"`
	Status          bool                `json:"status"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Media           []*MediaDto         `json:"media"`
}

type PostWithLikedDto struct {
	ID              uuid.UUID           `json:"id"`
	UserId          uuid.UUID           `json:"user_id"`
	User            *UserDto            `json:"user"`
	ParentId        *uuid.UUID          `json:"parent_id"`
	ParentPost      *PostDto            `json:"parent_post"`
	Content         string              `json:"content"`
	LikeCount       int                 `json:"like_count"`
	CommentCount    int                 `json:"comment_count"`
	Privacy         consts.PrivacyLevel `json:"privacy"`
	Location        string              `json:"location"`
	IsAdvertisement bool                `json:"is_advertisement"`
	Status          bool                `json:"status"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Media           []*MediaDto         `json:"media"`
	IsLiked         bool                `json:"is_liked"`
}

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

func ToPostDto(postResult common.PostResult) *PostDto {
	var parentPost *PostDto

	if postResult.ParentPost != nil {
		parentPost = &PostDto{
			ID:              postResult.ParentPost.ID,
			UserId:          postResult.ParentPost.UserId,
			User:            ToUserDto(postResult.ParentPost.User),
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

	return &PostDto{
		ID:              postResult.ID,
		UserId:          postResult.UserId,
		User:            ToUserDto(postResult.User),
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

func ToPostWithLikedDto(
	postResult common.PostResultWithLiked,
) *PostWithLikedDto {
	var parentPost *PostDto

	if postResult.ParentPost != nil {
		parentPost = &PostDto{
			ID:              postResult.ParentPost.ID,
			UserId:          postResult.ParentPost.UserId,
			User:            ToUserDto(postResult.ParentPost.User),
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

	return &PostWithLikedDto{
		ID:              postResult.ID,
		UserId:          postResult.UserId,
		User:            ToUserDto(postResult.User),
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
		IsLiked:         postResult.IsLiked,
	}
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
