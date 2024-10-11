package repository_implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
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
	media *model.Media,
) (*model.Media, error) {
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
) (*model.Media, error) {
	var media model.Media

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
	res := r.db.WithContext(ctx).Delete(&model.Media{}, mediaId)
	return res.Error
}

func (r *rMedia) GetMedia(ctx context.Context, query interface{}, args ...interface{}) (*model.Media, error) {
	media := &model.Media{}

	if res := r.db.WithContext(ctx).Model(media).Where(query, args...).First(query); res.Error != nil {
		return nil, res.Error
	}

	return media, nil
}

func (r *rMedia) GetManyMedia(ctx context.Context) ([]*model.Media, error) {
	var medias []*model.Media
	if err := r.db.WithContext(ctx).Find(&medias).Error; err != nil {
		return nil, err
	}

	return medias, nil
}
