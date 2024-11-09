package repo_impl

import (
	"context"
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	setting_mapper "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/mapper"
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
	settingId uint,
) (*user_entity.Setting, error) {
	var settingModel models.Setting
	if err := r.db.WithContext(ctx).
		First(&settingModel, settingId).
		Error; err != nil {
		return nil, err
	}

	return setting_mapper.FromSettingModel(&settingModel), nil
}

func (r *rSetting) CreateOne(
	ctx context.Context,
	settingEntity *user_entity.Setting,
) (*user_entity.Setting, error) {
	settingModel := setting_mapper.ToSettingModel(settingEntity)

	if err := r.db.WithContext(ctx).
		Create(settingModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, settingModel.ID)
}

func (r *rSetting) UpdateOne(
	ctx context.Context,
	settingId uint,
	settingUpdateEntity *user_entity.SettingUpdate,
) (*user_entity.Setting, error) {
	updates := map[string]interface{}{}

	if settingUpdateEntity.Language != nil {
		updates["language"] = *settingUpdateEntity.Language
	}

	if settingUpdateEntity.UpdatedAt != nil {
		updates["updated_at"] = *settingUpdateEntity.UpdatedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Setting{}).
		Where("id = ?", settingId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, settingId)
}

func (r *rSetting) DeleteOne(
	ctx context.Context,
	settingId uint,
) error {
	res := r.db.WithContext(ctx).
		Delete(&models.Setting{}, settingId)
	return res.Error
}

func (r *rSetting) GetSetting(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*user_entity.Setting, error) {
	var settingModel models.Setting

	if err := r.db.WithContext(ctx).
		Model(&settingModel).
		Where(query, args...).
		First(settingModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, settingModel.ID)
}
