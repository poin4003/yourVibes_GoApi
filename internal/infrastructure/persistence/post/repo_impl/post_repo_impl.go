package repo_impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
	"gorm.io/gorm"
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("Media").
		Preload("ParentPost.Media").
		First(&postModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
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

	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return errors.New("no fields to update")
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
	postFound := &models.Post{}

	return mapper.FromPostModel(postFound), r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check if post exists
		if err := tx.WithContext(ctx).
			Model(postFound).
			Where("id = ?", id).
			Select("id, user_id").
			First(postFound).
			Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Delete related reports
		result := tx.WithContext(ctx).
			Where("id IN (?)",
				tx.Model(&models.PostReport{}).
					Select("report_id").
					Where("reported_post_id = ?", id),
			).
			Delete(&models.Report{})

		if result.Error != nil {
			return response.NewServerFailedError(result.Error.Error())
		}

		// 3. Update postCount -1 in User table
		if err := tx.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", postFound.UserId).
			Update("post_count", gorm.Expr("post_count - ?", 1)).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// 4. Delete related comment
		commentResult := tx.WithContext(ctx).
			Model(&models.Comment{}).
			Where("post_id = ?", id).
			Delete(&models.Comment{})

		if commentResult.Error != nil {
			return response.NewServerFailedError(commentResult.Error.Error())
		}

		// 5. Delete related report (Comment report)
		if commentResult.RowsAffected > 0 {
			commentReportResult := tx.WithContext(ctx).
				Where("id IN (?)",
					tx.Model(&models.CommentReport{}).
						Select("report_id").
						Where("reported_comment_id IN (?)",
							tx.Model(&models.Comment{}).
								Select("id").
								Where("post_id = ?", id),
						),
				).
				Delete(&models.Report{})

			if commentReportResult.Error != nil {
				return response.NewServerFailedError(commentReportResult.Error.Error())
			}
		}

		// 6. Delete post
		if err := tx.WithContext(ctx).
			Delete(&models.Post{}, "id = ?", id).
			Error; err != nil {
			return err
		}
		return nil
	})
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("Media").
		Preload("ParentPost.Media").
		First(&postModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
		sortColumn := ""
		switch query.SortBy {
		case "id":
			sortColumn = "id"
		case "title":
			sortColumn = "title"
		case "content":
			sortColumn = "content"
		case "created_at":
			sortColumn = "created_at"
		case "location":
			sortColumn = "location"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
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
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("Media").
		Preload("ParentPost.Media").
		Find(&postModels).
		Error; err != nil {
		return nil, nil, err
	}

	var postIds []uuid.UUID
	for _, post := range postModels {
		postIds = append(postIds, post.ID)
	}

	var likedPostIds []uuid.UUID
	if err := r.db.Model(&models.LikeUserPost{}).
		Select("post_id").
		Where("user_id = ? AND post_id IN ?", authenticatedUserId, postIds).
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

func (r *rPost) CheckPostOwner(
	ctx context.Context,
	postId uuid.UUID,
	userId uuid.UUID,
) (bool, error) {
	var ownerId string
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Select("user_id").
		Where("id = ?", postId).
		First(&ownerId).
		Error; err != nil {
		return false, err
	}

	ownerUUID, err := uuid.Parse(ownerId)
	if err != nil {
		return false, fmt.Errorf("invalid UUID for user_id: %v", err)
	}

	if ownerUUID != userId {
		return false, nil
	}

	return true, nil
}

func (r *rPost) GetTotalPostCount(ctx context.Context) (int, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Count(&total).
		Error; err != nil {
		return 0, err
	}

	return int(total), nil
}
