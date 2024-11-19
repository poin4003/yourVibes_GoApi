package mapper

import (
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
)

func ToMediaModel(media *entities.Media) *models.Media {
	m := &models.Media{
		PostId:    media.PostId,
		MediaUrl:  media.MediaUrl,
		Status:    media.Status,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}
	m.ID = media.ID

	return m
}

func FromMediaModel(m *models.Media) *entities.Media {
	var media = &entities.Media{
		PostId:    m.PostId,
		MediaUrl:  m.MediaUrl,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	media.ID = m.ID

	return media
}
