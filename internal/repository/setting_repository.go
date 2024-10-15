package repository

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
)

type (
	ISettingRepository interface {
		CreateSetting(ctx context.Context, setting *model.Setting) (*model.Setting, error)
		UpdateSetting(ctx context.Context, settingId uint, updateData map[string]interface{}) (*model.Setting, error)
		DeleteSetting(ctx context.Context, settingId uint) error
		GetSetting(ctx context.Context, query interface{}, args ...interface{}) (*model.Setting, error)
	}
)

var (
	localSetting ISettingRepository
)

func Setting() ISettingRepository {
	if localSetting == nil {
		panic("repository_implement localSetting not found for interface ISetting")
	}

	return localSetting
}

func InitSettingRepository(i ISettingRepository) {
	localSetting = i
}
