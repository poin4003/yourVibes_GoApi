package repo_impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/advertise/mapper"
	"gorm.io/gorm"
)

type rAdvertise struct {
	db *gorm.DB
}

func NewAdvertiseRepositoryImplement(db *gorm.DB) *rAdvertise {
	return &rAdvertise{db: db}
}

func (r *rAdvertise) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Advertise, error) {
	var advertiseModel models.Advertise
	if err := r.db.WithContext(ctx).
		Preload("Bill").
		First(&advertiseModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.FromAdvertiseModel(&advertiseModel), nil
}

func (r *rAdvertise) GetOne(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Advertise, error) {
	var advertiseModel models.Advertise

	if err := r.db.WithContext(ctx).
		Model(&advertiseModel).
		Preload("Bill").
		Preload("Post.User").
		Preload("Post.Media").
		Preload("Post.ParentPost.Media").
		Preload("Post.ParentPost.User").
		First(&advertiseModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.FromAdvertiseModelForAdvertiseDetail(&advertiseModel), nil
}

func (r *rAdvertise) GetDetailAndStatisticOfAdvertise(
	ctx context.Context,
	id uuid.UUID,
) (*entities.AdvertiseForStatistic, error) {
	advertiseModel := &models.Advertise{}

	if err := r.db.WithContext(ctx).
		Model(advertiseModel).
		First(&advertiseModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	type StatsSummary struct {
		TotalReach       int
		TotalClicks      int
		TotalImpressions int
	}
	var statsSummary StatsSummary
	if err := r.db.WithContext(ctx).
		Model(&models.Statistics{}).
		Select(
			"COALESCE(SUM(reach), 0) as total_reach, "+
				"COALESCE(SUM(clicks), 0) as total_clicks, "+
				"COALESCE(SUM(impression), 0) as total_impressions",
		).
		Where(
			"post_id = ? AND created_at BETWEEN ? AND ?",
			advertiseModel.PostId,
			advertiseModel.StartDate,
			advertiseModel.EndDate,
		).
		Scan(&statsSummary).
		Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	statEntities := r.getAverageStatistics(
		ctx,
		advertiseModel.PostId,
		&advertiseModel.StartDate, &advertiseModel.EndDate,
		10,
	)

	return mapper.FromAdvertiseModelForDetailAndStatistics(
		advertiseModel,
		statsSummary.TotalReach,
		statsSummary.TotalClicks,
		statsSummary.TotalImpressions,
		statEntities,
	), nil
}

func (r *rAdvertise) GetMany(
	ctx context.Context,
	query *query.GetManyAdvertiseQuery,
) ([]*entities.Advertise, *response.PagingResponse, error) {
	var advertises []*models.Advertise
	var total int64

	db := r.db.WithContext(ctx).
		Model(&models.Advertise{}).
		Joins("JOIN posts ON posts.id = advertises.post_id").
		Joins("JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN bills ON bills.advertise_id = advertises.id")

	if query.PostId != uuid.Nil {
		db = db.Where("advertises.post_id = ?", query.PostId)
	}

	if query.UserEmail != "" {
		db = db.Where("users.email = ?", query.UserEmail)
	}

	if query.Status != nil {
		db = db.Where("bills.status = ?", query.Status)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("advertises.created_at >= ?", query.FromDate)
	}
	if !query.ToDate.IsZero() {
		db = db.Where("advertises.created_at <= ?", query.ToDate)
	}

	if query.FromPrice > 0 {
		db = db.Where("bills.price >= ?", query.FromPrice)
	}
	if query.ToPrice > 0 {
		db = db.Where("bills.price <= ?", query.ToPrice)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "price":
			sortColumn = "bills.price"
		case "created_at":
			sortColumn = "advertises.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
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

	if err := db.Count(&total).
		Offset(offset).
		Limit(limit).
		Preload("Bill").
		Preload("Post", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, user_id")
		}).
		Preload("Post.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("users.id, users.email")
		}).
		Find(&advertises).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var advertiseEntities []*entities.Advertise
	for _, advertise := range advertises {
		advertiseEntities = append(advertiseEntities, mapper.FromAdvertiseModel(advertise))
	}

	return advertiseEntities, pagingResponse, nil
}

func (r *rAdvertise) CreateOne(
	ctx context.Context,
	entity *entities.Advertise,
) (*entities.Advertise, error) {
	advertiseModel := mapper.ToAdvertiseModel(entity)

	if err := r.db.WithContext(ctx).
		Create(advertiseModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, advertiseModel.ID)
}

func (r *rAdvertise) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.AdvertiseUpdate,
) (*entities.Advertise, error) {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no field to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Advertise{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rAdvertise) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.Advertise{}, id).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *rAdvertise) GetLatestAdsByPostId(
	ctx context.Context,
	postId uuid.UUID,
) (*entities.Advertise, error) {
	var advertise models.Advertise

	if err := r.db.WithContext(ctx).
		Where("post_id = ?", postId).
		Order("created_at desc").
		First(&advertise).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.GetById(ctx, advertise.ID)
}

func (r *rAdvertise) CheckExists(
	ctx context.Context,
	postId uuid.UUID,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.Advertise{}).
		Where("post_id = ?", postId).
		Count(&count).
		Error; err != nil {
	}

	return count > 0, nil
}

func (r *rAdvertise) DeleteMany(
	ctx context.Context,
	condition map[string]interface{},
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Where("advertise_id IN (?)",
				tx.Model(&models.Advertise{}).
					Select("id").
					Where(condition)).
			Delete(&models.Bill{}).
			Error; err != nil {
			return fmt.Errorf("failed to delete bills: %w", err)
		}

		if err := tx.WithContext(ctx).
			Where(condition).
			Delete(&models.Advertise{}).
			Error; err != nil {
			return fmt.Errorf("failed to delete advertise: %w", err)
		}

		return nil
	})
}

func (r *rAdvertise) getAverageStatistics(
	ctx context.Context,
	postId uuid.UUID,
	startDate, endDate *time.Time,
	numIntervals int,
) []*entities.StatisticEntity {
	if startDate == nil || endDate == nil {
		return nil
	}

	duration := endDate.Sub(*startDate)
	if duration <= 0 {
		return nil
	}

	intervalDuration := duration / time.Duration(numIntervals)

	statEntities := make([]*entities.StatisticEntity, numIntervals)
	for i := 0; i < numIntervals; i++ {
		intervalStart := startDate.Add(time.Duration(i) * intervalDuration)
		intervalEnd := startDate.Add(time.Duration(i+1) * intervalDuration)
		if i == numIntervals-1 {
			intervalEnd = *endDate
		}

		type AvgStats struct {
			AvgReach       float64
			AvgClicks      float64
			AvgImpressions float64
		}
		var avgStats AvgStats
		if err := r.db.WithContext(ctx).
			Model(&models.Statistics{}).
			Select(
				"COALESCE(AVG(reach), 0) as avg_reach, "+
					"COALESCE(AVG(clicks), 0) as avg_clicks, "+
					"COALESCE(AVG(impression), 0) as avg_impressions",
			).
			Where("post_id = ? AND deleted_at IS NULL AND created_at BETWEEN ? AND ?", postId, intervalStart, intervalEnd).
			Scan(&avgStats).
			Error; err != nil {
			global.Logger.Error("Failed to calculate average statistic", zap.Error(err))
			continue
		}

		statEntities[i] = &entities.StatisticEntity{
			PostId:          postId,
			Reach:           int(avgStats.AvgReach),
			Clicks:          int(avgStats.AvgClicks),
			Impression:      int(avgStats.AvgImpressions),
			AggregationDate: intervalStart,
		}
	}

	return statEntities
}
