package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/query"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"gorm.io/gorm"
)

type rNewFeed struct {
	db *gorm.DB
}

func NewNewFeedRepositoryImplement(db *gorm.DB) *rNewFeed {
	return &rNewFeed{db: db}
}

func (r *rNewFeed) CreateMany(
	ctx context.Context,
	postId uuid.UUID,
	userId uuid.UUID,
) error {
	query := `
		INSERT INTO new_feeds (user_id, post_id)
		SELECT friend.friend_id, CAST(? AS UUID)
		FROM friends friend
		WHERE friend.user_id = CAST(? AS UUID)
		  AND NOT EXISTS (
			  SELECT 1
			  FROM new_feeds nf
			  WHERE nf.user_id = friend.friend_id
				AND nf.post_id = CAST(? AS UUID)
		  )
		UNION ALL
		SELECT CAST(? AS UUID), CAST(? AS UUID)
		WHERE NOT EXISTS (
			SELECT 1
			FROM new_feeds nf
			WHERE nf.user_id = CAST(? AS UUID)
			  AND nf.post_id = CAST(? AS UUID)
		);
	`

	if err := r.db.WithContext(ctx).
		Exec(query, postId, userId, postId, userId, postId, userId, postId).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *rNewFeed) DeleteOne(
	ctx context.Context,
	userId uuid.UUID,
	postId uuid.UUID,
) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userId, postId).
		Delete(&models.NewFeed{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *rNewFeed) DeleteMany(
	ctx context.Context,
	condition map[string]interface{},
) error {
	if err := r.db.WithContext(ctx).
		Model(models.NewFeed{}).
		Where(condition).
		Delete(&models.NewFeed{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rNewFeed) GetMany(
	ctx context.Context,
	query *query.GetNewFeedQuery,
) ([]*entities.Post, *response.PagingResponse, error) {
	var posts []*models.Post
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

	authenticatedUserId := query.UserId

	friendSubQuery := r.db.Model(&models.Friend{}).
		Select("friend_id").
		Where("user_id = ?", authenticatedUserId)

	db := r.db.WithContext(ctx)

	if err := db.Model(&models.Post{}).
		Joins("JOIN new_feeds ON new_feeds.post_id = posts.id").
		Where("status = true").
		Where("new_feeds.user_id = ?", authenticatedUserId).
		Where(`
			(posts.privacy = ? OR 
			(posts.privacy = ? AND (posts.user_id IN (?) OR posts.user_id = ?)) OR
			(posts.privacy = ? AND posts.user_id = ?))
		`, consts.PUBLIC, consts.FRIEND_ONLY, friendSubQuery, authenticatedUserId, consts.PRIVATE, authenticatedUserId,
		).
		Count(&total).Error; err != nil {
		return nil, nil, err
	}

	if err := db.Model(&models.Post{}).
		Joins("JOIN new_feeds ON new_feeds.post_id = posts.id").
		Where("status = true").
		Where("deleted_at IS NULL").
		Where("new_feeds.user_id = ?", query.UserId).
		Where(`
			(posts.privacy = ? OR 
			(posts.privacy = ? AND (posts.user_id IN (?) OR posts.user_id = ?)) OR
			(posts.privacy = ? AND posts.user_id = ?))
		`, consts.PUBLIC, consts.FRIEND_ONLY, friendSubQuery, authenticatedUserId, consts.PRIVATE, authenticatedUserId,
		).
		Select(`posts.*,
	   	EXISTS (
	       SELECT 1
	       FROM like_user_posts
	       WHERE like_user_posts.post_id = posts.id AND like_user_posts.user_id = ?
	   	) AS is_liked
		`, authenticatedUserId).
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
		Order("posts.updated_at desc").
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postEntities []*entities.Post
	for _, post := range posts {
		postEntity := mapper.FromPostModel(post)
		postEntities = append(postEntities, postEntity)
	}

	return postEntities, pagingResponse, nil
}

func (r *rNewFeed) CreateManyWithRandomUser(
	ctx context.Context,
	numUsers int,
) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
        WITH 
            inserted AS (
                INSERT INTO new_feeds (user_id, post_id)
                SELECT u.id, a.post_id
                FROM users u
                CROSS JOIN (
                    SELECT advertises.id AS advertise_id, advertises.post_id
                    FROM advertises
                    JOIN bills ON bills.advertise_id = advertises.id
                    WHERE bills.status = true
                    AND advertises.start_date <= ?
                    AND advertises.end_date >= ?
                    AND advertises.deleted_at IS NULL
                    AND bills.deleted_at IS NULL
                ) a
                WHERE NOT EXISTS (
                    SELECT 1
                    FROM new_feeds nf
                    WHERE nf.user_id = u.id
                    AND nf.post_id = a.post_id
                )
                ORDER BY RANDOM()
                LIMIT ?
                RETURNING post_id
            ),
            reach_counts AS (
                SELECT post_id, COUNT(*) as reach_count
                FROM inserted
                GROUP BY post_id
            )
        UPDATE statistics
        SET reach = statistics.reach + rc.reach_count,
            updated_at = ?
        FROM reach_counts rc
        WHERE statistics.post_id = rc.post_id;
    `
	now := time.Now().Truncate(24 * time.Hour)

	result := tx.Exec(query, now, now, numUsers, now)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *rNewFeed) DeleteExpiredAdvertiseFromNewFeeds(
	ctx context.Context,
) error {
	query := `
        DELETE FROM new_feeds
        WHERE post_id IN (
            SELECT posts.id
            FROM posts
            WHERE EXISTS (
                SELECT 1
                FROM advertises
                WHERE advertises.post_id = posts.id
                  AND advertises.end_date < ?
            )
        )
    `

	now := time.Now()
	endOfToday := now.Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)

	if err := r.db.WithContext(ctx).
		Exec(query, endOfToday).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rNewFeed) CreateManyFeaturedPosts(
	ctx context.Context,
	numUsers int,
) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	now := time.Now()
	averageTimeToPush := now.AddDate(0, 0, -7)

	query := `
       	WITH
			latest_statistics AS (
				SELECT DISTINCT ON (post_id) post_id, clicks, impression
				FROM statistics
				WHERE deleted_at IS NULL
				ORDER BY post_id, created_at DESC
			),
			featured_posts AS (
				SELECT p.id AS post_id
				FROM posts p
				JOIN latest_statistics ls ON ls.post_id = p.id
				WHERE p.like_count >= ?
				  AND p.comment_count >= ?
				  AND ls.clicks >= ?
				  AND ls.impression >= ?
				  AND p.privacy = 'public'
				  AND p.is_advertisement = 0
				  AND p.status = true
				  AND p.created_at >= ?
				  AND p.created_at <= ?
				  AND p.deleted_at IS NULL
			),
			inserted AS (
				INSERT INTO new_feeds (user_id, post_id)
				SELECT u.id, fp.post_id
				FROM users u
				CROSS JOIN featured_posts fp
				WHERE NOT EXISTS (
					SELECT 1
					FROM new_feeds nf
					WHERE nf.user_id = u.id
					  AND nf.post_id = fp.post_id
				)
				ORDER BY RANDOM()
				LIMIT ?
				RETURNING post_id
			),
			reach_counts AS (
				SELECT post_id, COUNT(*) as reach_count
				FROM inserted
				GROUP BY post_id
			),
			update_stats AS (
				UPDATE statistics
				SET reach = statistics.reach + rc.reach_count,
					updated_at = ?
				FROM reach_counts rc
				WHERE statistics.post_id = rc.post_id
				  AND statistics.created_at = (
					  SELECT MAX(created_at)
					  FROM statistics s
					  WHERE s.post_id = rc.post_id
						AND s.deleted_at IS NULL
				  )
				RETURNING rc.post_id
			)
		UPDATE posts
		SET is_featured_post = true,
			updated_at = ?
		WHERE id IN (SELECT post_id FROM inserted);
    `

	result := tx.Exec(query, 3, 5, 10, 10, averageTimeToPush, now, numUsers, now, now)
	if result.Error != nil {
		return result.Error
	}

	global.Logger.Info("Pushed featured posts to newfeed", zap.Int("num_users", numUsers), zap.Int64("rows_affected", result.RowsAffected))

	return nil
}

func (r *rNewFeed) DeleteExpiredFeaturedPostsFromNewFeeds(ctx context.Context) error {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	query := `
		WITH expired_posts AS (
			SELECT id
			FROM posts
			WHERE is_featured_post = true
				AND created_at < ?
				AND deleted_at IS NULL
		),
		deleted_new_feeds AS (
			DELETE FROM new_feeds
			WHERE post_id IN (SELECT id FROM expired_posts)
			RETURNING post_id
		)
		UPDATE posts
		SET is_featured_post = false
		WHERE id IN (SELECT post_id FROM deleted_new_feeds)
	`

	if err := r.db.WithContext(ctx).Exec(query, sevenDaysAgo).Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rNewFeed) ExpireAdvertiseByPostID(ctx context.Context, postID uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Advertise{}).
			Where("post_id = ?", postID).
			Update("end_date", yesterday).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Post{}).
			Where("id = ?", postID).
			Update("is_advertisement", 2).Error; err != nil {
			return err
		}

		return nil
	})
}
