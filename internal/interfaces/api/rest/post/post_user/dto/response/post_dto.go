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
	ParentPost      *PostWithLikedDto   `json:"parent_post"`
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
	ParentPost      *PostWithLikedDto   `json:"parent_post"`
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

func ToPostDto(postResult common.PostResult) *PostDto {
	var parentPost *PostWithLikedDto

	if postResult.ParentPost != nil {
		parentPost = &PostWithLikedDto{
			ID:              postResult.ID,
			UserId:          postResult.UserId,
			User:            ToUserDto(postResult.User),
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
	var parentPost *PostWithLikedDto

	if postResult.ParentPost != nil {
		parentPost = &PostWithLikedDto{
			ID:              postResult.ID,
			UserId:          postResult.UserId,
			User:            ToUserDto(postResult.User),
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
