package repo_impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"

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
		Where("posts.id = ? AND status = true", id).
		Where("deleted_at IS NULL").
		Preload("Media").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", true)
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.Media").
		First(&postModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return mapper.FromPostModel(&postModel), nil
}

func (r *rPost) CreateOne(
	ctx context.Context,
	entity *entities.Post,
) (*entities.Post, error) {
	postModel := mapper.ToPostModel(entity)

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Create(postModel).
			Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Create(&models.Statistics{
				ID:     uuid.New(),
				PostId: postModel.ID,
			}).
			Error; err != nil {
			return err
		}

		if err := r.updatePostCountByUserId(ctx, tx, entity.UserId, 1); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return r.GetById(ctx, postModel.ID)
}

func (r *rPost) updatePostCountByUserId(
	ctx context.Context,
	tx *gorm.DB,
	id uuid.UUID,
	countChange int,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn(
			"post_count", gorm.Expr("GREATEST(post_count + ?, 0)", countChange),
		).Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
		if err := r.updatePostCountByUserId(ctx, tx, postFound.UserId, -1); err != nil {
			return err
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

		// 6. Delete statistic
		if err := tx.WithContext(ctx).
			Delete(&models.Statistics{}, "post_id = ?", id).
			Error; err != nil {
			return err
		}

		// 7. Delete post
		if err := tx.WithContext(ctx).
			Delete(&models.Post{}, "id = ?", id).
			Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *rPost) GetMany(
	ctx context.Context,
	query *query.GetManyPostQuery,
) ([]*entities.Post, *response.PagingResponse, error) {
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

	if query.IsAdvertisement != nil {
		db = db.Where("is_advertisement = ?", query.IsAdvertisement)
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
				db = db.Order(sortColumn)
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

	if err = db.Offset(offset).Limit(limit).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("Media").
		Preload("ParentPost", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", true)
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.Media").
		Find(&postModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postEntities []*entities.Post
	for _, post := range postModels {
		postEntity := mapper.FromPostModel(post)
		postEntities = append(postEntities, postEntity)
	}

	return postEntities, pagingResponse, nil
}

func (r *rPost) GetTrendingPost(
	ctx context.Context,
	query *query.GetTrendingPostQuery,
) ([]*entities.Post, *response.PagingResponse, error) {
	var postModels []*models.Post
	var total int64

	// 1. Paging
	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	now := time.Now()
	averageTimeToGet := now.AddDate(0, 0, -7)

	// 2. Get total record
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("status = ?", true).
		Where("deleted_at is NULL").
		Where("is_advertisement = ?", 0).
		Where("privacy = ?", consts.PUBLIC).
		Count(&total).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	// 3. subquery/CTE to calculate score
	subQuery := r.db.WithContext(ctx).
		Model(&models.Statistics{}).
		Select(`
			post_id, 
			COALESCE(SUM(reach), 0) AS total_reach, 
			COALESCE(SUM(clicks), 0) AS total_clicks, 
			COALESCE(SUM(impression), 0) AS total_impression`,
		).
		Group("post_id")

	if err := r.db.WithContext(ctx).
		Unscoped().
		Table("(?) AS s", subQuery).
		Select(
			"p.*, "+
				"(COALESCE(s.total_impression, 0) * 0.3 + "+
				"p.like_count * 0.25 + "+
				"p.comment_count * 0.2 + "+
				"COALESCE(s.total_clicks, 0) * 0.15 + "+
				"COALESCE(s.total_reach, 0) * 0.1) AS score",
		).
		Joins("RIGHT JOIN posts p ON p.id = s.post_id").
		Where("p.status = ?", true).
		Where("p.is_advertisement = ?", 0).
		Where("p.privacy = ?", consts.PUBLIC).
		Where("p.created_at >= ? AND p.created_at <= ?", averageTimeToGet, now).
		Where("deleted_at is NULL").
		Order("score DESC").
		Offset(offset).
		Limit(limit).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("Media").
		Preload("ParentPost", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", true)
		}).
		Preload("ParentPost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentPost.Media").
		Find(&postModels).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	// 5. Map to entity
	var postEntities []*entities.Post
	for _, post := range postModels {
		postEntity := mapper.FromPostModel(post)
		postEntities = append(postEntities, postEntity)
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return postEntities, pagingResponse, nil
}

func (r *rPost) UpdateExpiredAdvertisements(
	ctx context.Context,
) error {
	query := `
      	UPDATE posts
		SET is_advertisement = 2
		WHERE is_advertisement = 1
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

func (r *rPost) GetTotalPostCountByUserId(
	ctx context.Context,
	userId uuid.UUID,
) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("user_id = ?", userId).
		Count(&total).
		Error; err != nil {
		return 0, response.NewServerFailedError(err.Error())
	}

	return total, nil
}
