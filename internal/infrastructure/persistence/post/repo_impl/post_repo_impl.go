package repo_impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
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

func (r *rPost) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Post, error) {
	var postModel models.Post
	if err := r.db.WithContext(ctx).
		Preload("Media").
		Preload("User").
		Preload("ParentPost.User").
		Preload("ParentPost.Media").
		First(&postModel, id).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromPostModel(&postModel), nil
}

func (r *rPost) CreateOne(
	ctx context.Context,
	entity *entities.Post,
) (*entities.Post, error) {
	postModel := mapper.ToPostModel(entity)

	if err := r.db.WithContext(ctx).
		Create(postModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, postModel.ID)
}

func (r *rPost) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.PostUpdate,
) (*entities.Post, error) {
	updates := map[string]interface{}{}

	if updateData.Content != nil {
		updates["content"] = *updateData.Content
	}

	if updateData.LikeCount != nil {
		updates["like_count"] = *updateData.LikeCount
	}

	if updateData.CommentCount != nil {
		updates["comment_count"] = *updateData.CommentCount
	}

	if updateData.Privacy != nil {
		updates["privacy"] = *updateData.Privacy
	}

	if updateData.Location != nil {
		updates["location"] = *updateData.Location
	}

	if updateData.IsAdvertisement != nil {
		updates["is_advertisement"] = *updateData.IsAdvertisement
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rPost) UpdateMany(
	ctx context.Context,
	condition map[string]interface{},
	updateData *entities.PostUpdate,
) error {
	if len(condition) == 0 {
		return errors.New("condition cannot be empty")
	}

	updates := map[string]interface{}{}

	if updateData.Content != nil {
		updates["content"] = *updateData.Content
	}

	if updateData.LikeCount != nil {
		updates["like_count"] = *updateData.LikeCount
	}

	if updateData.CommentCount != nil {
		updates["comment_count"] = *updateData.CommentCount
	}

	if updateData.Privacy != nil {
		updates["privacy"] = *updateData.Privacy
	}

	if updateData.Location != nil {
		updates["location"] = *updateData.Location
	}

	if updateData.IsAdvertisement != nil {
		updates["is_advertisement"] = *updateData.IsAdvertisement
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
	}

	if len(updates) == 0 {
		return errors.New("no updates specified")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where(condition).
		Updates(updates).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rPost) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Post, error) {
	post, err := r.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := r.db.WithContext(ctx).Delete(mapper.ToPostModel(post))
	if res.Error != nil {
		return nil, res.Error
	}

	return post, nil
}

func (r *rPost) GetOne(
	ctx context.Context,
	id uuid.UUID,
	authenticatedUserId uuid.UUID,
) (*entities.PostWithLiked, error) {
	var postModel models.Post

	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("posts.id = ? AND status = true", id).
		Preload("Media").
		Preload("User").
		Preload("ParentPost.User").
		Preload("ParentPost.Media").
		First(&postModel).
		Error; err != nil {
		return nil, err
	}

	var isLiked bool
	if err := r.db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM like_user_posts
			WHERE like_user_posts.post_id = ? AND like_user_posts.user_id = ?
		)`, id, authenticatedUserId).
		Scan(&isLiked).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromPostWithLikedModel(&postModel, isLiked), nil
}

func (r *rPost) GetMany(
	ctx context.Context,
	query *query.GetManyPostQuery,
) ([]*entities.PostWithLiked, *response.PagingResponse, error) {
	var postModels []*models.Post
	var total int64

	authenticatedUserId := query.AuthenticatedUserId
	friendSubQuery := r.db.Model(&models.Friend{}).
		Select("friend_id").
		Where("user_id = ?", authenticatedUserId)

	db := r.db.WithContext(ctx).Model(&models.Post{})

	db = db.Where(
		"(privacy = ? OR (privacy = ? AND (user_id IN (?) OR user_id = ?)) OR (privacy = ? AND user_id = ?))",
		consts.PUBLIC, consts.FRIEND_ONLY, friendSubQuery, authenticatedUserId, consts.PRIVATE, authenticatedUserId,
	)

	db = db.Where("status = ?", true)

	if query.UserID != uuid.Nil {
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

	if err := db.Offset(offset).Limit(limit).
		Preload("User").
		Preload("Media").
		Preload("ParentPost.User").
		Preload("ParentPost.Media").
		Find(&postModels).
		Error; err != nil {
		return nil, nil, err
	}

	var likedPostIds []uuid.UUID
	if err := r.db.Model(&models.LikeUserPost{}).
		Select("post_id").
		Where("user_id = ?", authenticatedUserId).
		Find(&likedPostIds).
		Error; err != nil {
		return nil, nil, err
	}

	likedMap := make(map[uuid.UUID]bool)
	for _, id := range likedPostIds {
		likedMap[id] = true
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postEntities []*entities.PostWithLiked
	for _, post := range postModels {
		postEntity := mapper.FromPostWithLikedModel(post, likedMap[post.ID])
		postEntities = append(postEntities, postEntity)
	}

	return postEntities, pagingResponse, nil
}

func (r *rPost) UpdateExpiredAdvertisements(
	ctx context.Context,
) error {
	query := `
      	UPDATE posts
		SET is_advertisement = false
		WHERE is_advertisement = true 
		  AND NOT EXISTS (
			  SELECT 1
			  FROM advertises 
			  WHERE advertises.post_id = posts.id
				AND advertises.end_date >= ?
		  )
		  AND EXISTS (
			  SELECT 1
			  FROM advertises 
			  WHERE advertises.post_id = posts.id
				AND advertises.end_date < ?
		  ) 
    `

	now := time.Now()
	endOfToday := now.Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)

	if err := r.db.WithContext(ctx).
		Exec(query, endOfToday, endOfToday).
		Error; err != nil {
		return err
	}

	return nil
}
