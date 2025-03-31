package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
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
	entity *entities.LikeUserPost,
) error {
	if err := r.db.WithContext(ctx).
		Create(mapper.ToLikeUserPostModel(entity)).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rLikeUserPost) DeleteLikeUserPost(
	ctx context.Context,
	entity *entities.LikeUserPost,
) error {
	if err := r.db.WithContext(ctx).
		Delete(mapper.ToLikeUserPostModel(entity)).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rLikeUserPost) GetLikeUserPost(
	ctx context.Context,
	query *query.GetPostLikeQuery,
) ([]*entities.User, *response.PagingResponse, error) {
	var users []*models.User
	var total int64

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&models.User{})

	err := db.Joins("JOIN like_user_posts ON like_user_posts.user_id = users.id").
		Select("id, family_name, name, avatar_url").
		Where("like_user_posts.post_id = ?", query.PostId).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var userEntities []*entities.User
	for _, user := range users {
		userEntity := mapper.FromUserModel(user)
		userEntities = append(userEntities, userEntity)
	}

	return userEntities, pagingResponse, nil
}

func (r *rLikeUserPost) CheckUserLikePost(
	ctx context.Context,
	entity *entities.LikeUserPost,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.LikeUserPost{}).
		Where("post_id = ? AND user_id =?", entity.PostId, entity.UserId).
		Count(&count).Error; err != nil {
	}
	return count > 0, nil
}

func (r *rLikeUserPost) CheckUserLikeManyPost(
	ctx context.Context,
	query *query.CheckUserLikeManyPostQuery,
) (map[uuid.UUID]bool, error) {
	var likedPostsIds []uuid.UUID
	if err := r.db.WithContext(ctx).Model(&models.LikeUserPost{}).
		Select("post_id").
		Where("user_id = ? AND post_id IN ?", query.AuthenticatedUserId, query.PostIds).
		Find(&likedPostsIds).
		Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	likedMap := make(map[uuid.UUID]bool)
	for _, id := range likedPostsIds {
		likedMap[id] = true
	}

	return likedMap, nil
}
