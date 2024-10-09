package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

type (
	IPostRepository interface {
		CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}) (*model.Post, error)
		DeletePost(ctx context.Context, postId uuid.UUID) error
		GetPost(ctx context.Context, query interface{}, args ...interface{}) (*model.Post, error)
		GetAllPost(ctx context.Context) ([]*model.Post, error)
	}
)

var (
	localPost IPostRepository
)

func Post() IPostRepository {
	if localPost == nil {
		panic("repository_implement localPost not found for interface IPost")
	}

	return localPost
}

func InitPostRepository(i IPostRepository) {
	localPost = i
}
