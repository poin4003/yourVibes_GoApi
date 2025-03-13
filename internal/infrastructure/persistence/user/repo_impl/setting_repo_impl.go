package repo_impl

import (
	"context"
	"errors"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"

	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
	"gorm.io/gorm"
)

type rSetting struct {
	db *gorm.DB
}

func NewSettingRepositoryImplement(db *gorm.DB) *rSetting {
	return &rSetting{db: db}
}

func (r *rSetting) GetById(
	ctx context.Context,
	id uint,
) (*entities.Setting, error) {
	var settingModel models.Setting
	if err := r.db.WithContext(ctx).
		First(&settingModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.FromSettingModel(&settingModel), nil
}

func (r *rSetting) CreateOne(
	ctx context.Context,
	entity *entities.Setting,
) (*entities.Setting, error) {
	settingModel := mapper.ToSettingModel(entity)

	if err := r.db.WithContext(ctx).
		Create(settingModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, settingModel.ID)
}

func (r *rSetting) UpdateOne(
	ctx context.Context,
	id uint,
	updateData *entities.SettingUpdate,
) (*entities.Setting, error) {
	updates := converter.StructToMap(updateData)

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Setting{}).
		Where("id = ?", id).
		Updates(&updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rSetting) DeleteOne(
	ctx context.Context,
	id uint,
) error {
	res := r.db.WithContext(ctx).
		Delete(&models.Setting{}, id)
	return res.Error
}

func (r *rSetting) GetSetting(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Setting, error) {
	var settingModel models.Setting

	if err := r.db.WithContext(ctx).
		Model(&settingModel).
		Where(query, args).
		First(&settingModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.GetById(ctx, settingModel.ID)
}
