package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"gorm.io/gorm"
)

type rPost struct {
	db *gorm.DB
}

func NewPostRepositoryImplement(db *gorm.DB) *rPost {
	return &rPost{db: db}
}

func (r *rPost) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	return &model.Post{}, nil
}

func (r *rPost) UpdatePost(
	ctx context.Context,
	postId uuid.UUID,
	updateData map[string]interface{},
) (*model.Post, error) {
	return &model.Post{}, nil
}

func (r *rPost) DeletePost(ctx context.Context, postId uuid.UUID) error {
	return nil
}

func (r *rPost) GetPost(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Post, error) {
	return &model.Post{}, nil
}

func (r *rPost) GetAllPost(ctx context.Context) ([]*model.Post, error) {
	return []*model.Post{}, nil
}
