package repo_impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/statistic/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/statistic/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rStatistic struct {
	db *gorm.DB
}

func NewStatisticRepository(db *gorm.DB) *rStatistic {
	return &rStatistic{db: db}
}

func (r *rStatistic) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.StatisticEntity, error) {
	var statisticModel models.Statistics
	if err := r.db.WithContext(ctx).
		First(&statisticModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return mapper.FromStatisticModel(&statisticModel), nil
}

func (r *rStatistic) UpsertLatestStatistic(
	ctx context.Context,
	postId uuid.UUID,
	entity *entities.StatisticEntity,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Find post
		postFound, err := r.getPost(ctx, tx, postId)
		if err != nil {
			return err
		}

		statisticModel := &models.Statistics{}
		switch postFound.IsAdvertisement {
		// Create or update into day
		case consts.NOT_ADVERTISE, consts.WAS_ADVERTISE:
			statisticModel, err = r.getLatestStatistic(ctx, tx, postId)
			if err != nil {
				return err
			}
			if statisticModel == nil {
				statisticModel = mapper.ToStatisticModel(entity)
				if err = r.createOne(ctx, tx, entity); err != nil {
					return err
				}
			}
			if err = r.increaseParameter(
				ctx, tx, statisticModel.ID,
				entity.Reach, entity.Clicks, entity.Impression); err != nil {
				return err
			}
		// Create or update by day
		case consts.IS_ADVERTISE:
			statisticModel, err = r.getTodayStatistic(ctx, tx, postId)
			if err != nil {
				return err
			}
			if statisticModel == nil {
				statisticModel = mapper.ToStatisticModel(entity)
				if err = r.createOne(ctx, tx, entity); err != nil {
					return err
				}
			}
			if err = r.increaseParameter(
				ctx, tx, statisticModel.ID,
				entity.Reach, entity.Clicks, entity.Impression); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *rStatistic) createOne(
	ctx context.Context,
	tx *gorm.DB,
	entity *entities.StatisticEntity,
) error {
	statisticModel := mapper.ToStatisticModel(entity)

	if err := tx.WithContext(ctx).
		Create(statisticModel).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rStatistic) getPost(
	ctx context.Context,
	tx *gorm.DB,
	postId uuid.UUID,
) (*models.Post, error) {
	postFound := &models.Post{}
	if err := tx.WithContext(ctx).
		Model(postFound).
		Select("id, is_advertisement").
		First(postFound, postId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return postFound, nil
}

func (r *rStatistic) getLatestStatistic(
	ctx context.Context,
	tx *gorm.DB,
	postId uuid.UUID,
) (*models.Statistics, error) {
	statisticModel := &models.Statistics{}
	if err := tx.WithContext(ctx).
		Model(statisticModel).
		Where("post_id = ?", postId).
		Order("created_at desc").
		Select("id").
		First(statisticModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return statisticModel, nil
}

func (r *rStatistic) getTodayStatistic(
	ctx context.Context,
	tx *gorm.DB,
	postId uuid.UUID,
) (*models.Statistics, error) {
	now := time.Now()
	startOfDay := now.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	statisticModel := &models.Statistics{}
	if err := tx.WithContext(ctx).
		Model(statisticModel).
		Where("post_id = ? AND created_at >= ? AND created_at < ?", postId, startOfDay, endOfDay).
		Select("id").
		First(statisticModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return statisticModel, nil
}

func (r *rStatistic) increaseParameter(
	ctx context.Context,
	tx *gorm.DB,
	id uuid.UUID,
	reach, clicks, impression int,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Statistics{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reach":      gorm.Expr("reach + ?", reach),
			"clicks":     gorm.Expr("clicks + ?", clicks),
			"impression": gorm.Expr("impression + ?", impression),
		}).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}
