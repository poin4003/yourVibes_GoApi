package repo_impl

import (
	"context"
	"github.com/google/uuid"
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
		return nil, err
	}

	return mapper.FromAdvertiseModel(&advertiseModel), nil
}

func (r *rAdvertise) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Advertise, error) {
	var advertiseModel models.Advertise

	if err := r.db.WithContext(ctx).
		Model(&advertiseModel).
		Where(query, args...).
		First(&advertiseModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, advertiseModel.ID)
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
	updates := map[string]interface{}{}

	if updateData.StartDate != nil {
		updates["start_date"] = *updateData.StartDate
	}

	if updateData.EndDate != nil {
		updates["end_date"] = *updateData.EndDate
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
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
