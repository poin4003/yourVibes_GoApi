package repository_implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"gorm.io/gorm"
)

type rSetting struct {
	db *gorm.DB
}

func NewSettingRepositoryImplement(db *gorm.DB) *rSetting {
	return &rSetting{db: db}
}

func (r *rSetting) CreateSetting(
	ctx context.Context,
	setting *model.Setting,
) (*model.Setting, error) {
	res := r.db.WithContext(ctx).Create(setting)

	if res.Error != nil {
		return nil, res.Error
	}

	return setting, nil
}

func (r *rSetting) UpdateSetting(
	ctx context.Context,
	settingId uint,
	updateData map[string]interface{},
) (*model.Setting, error) {
	var setting model.Setting

	if err := r.db.WithContext(ctx).First(&setting, settingId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&setting).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

func (r *rSetting) DeleteSetting(
	ctx context.Context,
	settingId uint,
) error {
	res := r.db.WithContext(ctx).Delete(&model.Setting{}, settingId)
	return res.Error
}

func (r *rSetting) GetSetting(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*model.Setting, error) {
	setting := &model.Setting{}

	if res := r.db.WithContext(ctx).Model(setting).Where(query, args...).First(setting); res.Error != nil {
		return nil, res.Error
	}

	return setting, nil
}
