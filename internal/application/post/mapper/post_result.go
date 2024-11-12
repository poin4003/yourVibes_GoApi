package mapper

import (
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
)

func NewPostWithIsLiked(
	post *post_entity.Post,
	isLike bool,
) *common.PostResult {
	if post == nil {
		return nil
	}

	parentId := &uuid.Nil
	if post.ParentId != nil {
		parentId = post.ParentId
	}

	return &common.PostResult{
		ID:              post.ID,
		UserId:          post.UserId,
		ParentId:        parentId,
		ParentPost:      nil,
		Content:         post.Content,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		Privacy:         post.Privacy,
		Location:        post.Location,
		IsAdvertisement: post.IsAdvertisement,
		Status:          post.Status,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		IsLiked:         isLike,
		Media:           []*common.MediaResult{},
	}
}
