package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

type rMedia struct {
	db *gorm.DB
}

func NewMediaRepositoryImplement(db *gorm.DB) *rMedia {
	return &rMedia{db: db}
}

func (r *rMedia) CreateMedia(
	ctx context.Context,
	media *models.Media,
) (*models.Media, error) {
	res := r.db.WithContext(ctx).Create(media)

	if res.Error != nil {
		return nil, res.Error
	}

	return media, nil
}

func (r *rMedia) UpdateMedia(
	ctx context.Context,
	mediaId uint,
	updateData map[string]interface{},
) (*models.Media, error) {
	var media models.Media

	if err := r.db.WithContext(ctx).First(&media, mediaId).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&media).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &media, nil
}

func (r *rMedia) DeleteMedia(
	ctx context.Context,
	mediaId uint,
) error {
	res := r.db.WithContext(ctx).Delete(&models.Media{}, mediaId)
	return res.Error
}

func (r *rMedia) GetMedia(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*models.Media, error) {
	media := &models.Media{}

	if res := r.db.WithContext(ctx).Model(media).Where(query, args...).First(media); res.Error != nil {
		return nil, res.Error
	}

	return media, nil
}

func (r *rMedia) GetManyMedia(
	ctx context.Context,
	query interface{},
	args ...interface{},
) ([]*models.Media, error) {
	var medias []*models.Media
	if err := r.db.WithContext(ctx).Where(query, args...).Find(&medias).Error; err != nil {
		return nil, err
	}

	return medias, nil
}
