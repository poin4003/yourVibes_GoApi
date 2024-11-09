package mapper

import (
	user_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/user/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToSettingModel(setting *user_entity.Setting) *models.Setting {
	s := &models.Setting{
		UserId:    setting.UserId,
		Language:  setting.Language,
		Status:    setting.Status,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
	s.ID = setting.ID

	return s
}

func FromSettingModel(s *models.Setting) *user_entity.Setting {
	var setting = &user_entity.Setting{
		UserId:    s.UserId,
		Language:  s.Language,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	setting.ID = s.ID

	return setting
}
