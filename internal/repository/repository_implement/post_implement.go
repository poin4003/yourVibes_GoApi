package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
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

	if err := r.db.WithContext(ctx).Preload("Media").Preload("User").First(&post, postId).Error; err != nil {
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
) (*model.Post, error) {
	post := &model.Post{}
	res := r.db.WithContext(ctx).First(post, postId)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Delete(post)
	if res.Error != nil {
		return nil, res.Error
	}

	return post, nil
}

func (r *rPost) GetPost(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Post, error) {
	post := &model.Post{}

	if res := r.db.WithContext(ctx).Model(post).Preload("Media").Preload("User").Where(query, args...).First(post); res.Error != nil {
		return nil, res.Error
	}

	return post, nil
}

func (r *rPost) GetManyPost(
	ctx context.Context,
	query *query_object.PostQueryObject,
) ([]*model.Post, *response.PagingResponse, error) {
	var posts []*model.Post
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Post{})

	if query.UserID != "" {
		db = db.Where("user_id = ?", query.UserID)
	}

	if query.Content != "" {
		db = db.Where("LOWER(content) LIKE LOWER(?)", "%"+query.Content+"%")
	}

	if !query.CreatedAt.IsZero() {
		createAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createAt)
	}

	if query.Location != "" {
		db = db.Where("LOWER(location) LIKE LOWER(?)", "%"+query.Location+"%")
	}

	if query.SortBy != "" {
		switch query.SortBy {
		case "id":
			if query.IsDescending {
				db = db.Order("id DESC")
			} else {
				db = db.Order("id ASC")
			}
		case "title":
			if query.IsDescending {
				db = db.Order("title DESC")
			} else {
				db = db.Order("title ASC")
			}
		case "content":
			if query.IsDescending {
				db = db.Order("content DESC")
			} else {
				db = db.Order("content ASC")
			}
		case "created_at":
			if query.IsDescending {
				db = db.Order("created_at DESC")
			} else {
				db = db.Order("created_at ASC")
			}
		case "location":
			if query.IsDescending {
				db = db.Order("location DESC")
			} else {
				db = db.Order("location ASC")
			}
		}
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Preload("Media").Preload("User").Find(&posts).Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return posts, pagingResponse, nil
}
