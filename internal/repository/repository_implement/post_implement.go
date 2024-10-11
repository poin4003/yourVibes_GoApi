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
	res := r.db.WithContext(ctx).Create(post)

	if res.Error != nil {
		return nil, res.Error
	}

	return post, nil
}

func (r *rPost) UpdatePost(
	ctx context.Context,
	postId uuid.UUID,
	updateData map[string]interface{},
) (*model.Post, error) {
	var post model.Post

	if err := r.db.WithContext(ctx).First(&post, postId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&post).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *rPost) DeletePost(
	ctx context.Context,
	postId uuid.UUID,
) error {
	res := r.db.WithContext(ctx).Delete(&model.Post{}, postId)
	return res.Error
}

func (r *rPost) GetPost(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Post, error) {
	post := &model.Post{}

	if res := r.db.WithContext(ctx).Model(post).Where(query, args...).First(post); res.Error != nil {
		return nil, res.Error
	}

	return post, nil
}

func (r *rPost) GetManyPost(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.db.WithContext(ctx).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
