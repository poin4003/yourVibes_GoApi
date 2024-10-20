package repository_implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"gorm.io/gorm"
)

type rLikeUserPost struct {
	db *gorm.DB
}

func NewLikeUserPostRepositoryImplement(db *gorm.DB) *rLikeUserPost {
	return &rLikeUserPost{db: db}
}

func (r *rLikeUserPost) CreateLikeUserPost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	res := r.db.WithContext(ctx).Create(likeUserPost)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rLikeUserPost) DeleteLikeUserPost(
	ctx context.Context,
	likeUserPost *model.LikeUserPost,
) error {
	res := r.db.WithContext(ctx).Delete(likeUserPost)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rLikeUserPost) GetLikeUserPost(
	ctx context.Context,
	query *query_object.PostLikeQueryObject,
) ([]*model.User, error) {
	var users []*model.User

	db := r.db.WithContext(ctx).Model(&model.User{})

	// Thay đổi mô hình để truy vấn từ bảng likes
	if query.PostID != "" {
		err := r.db.WithContext(ctx).
			Model(&model.User{}).
			Joins("JOIN like_user_posts ON like_user_posts.user_id = users.id").
			Where("like_user_posts.post_id = ?", query.PostID).
			Find(&users).Error
		if err != nil {
			return nil, err
		}
	} else {
		// Nếu không có postID, trả về danh sách tất cả người dùng hoặc xử lý khác
		err := r.db.WithContext(ctx).Find(&users).Error
		if err != nil {
			return nil, err
		}
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
	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
