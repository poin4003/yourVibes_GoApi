package repo_impl

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/mapper"
	"gorm.io/gorm"
)

type rMedia struct {
	db *gorm.DB
}

func NewMediaRepositoryImplement(db *gorm.DB) *rMedia {
	return &rMedia{db: db}
}

func (r *rMedia) GetById(
	ctx context.Context,
	id uint,
) (*entities.Media, error) {
	var mediaModel models.Media
	if err := r.db.WithContext(ctx).
		First(&mediaModel, id).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromMediaModel(&mediaModel), nil
}

func (r *rMedia) CreateOne(
	ctx context.Context,
	entity *entities.Media,
) (*entities.Media, error) {
	mediaModel := mapper.ToMediaModel(entity)

	if err := r.db.WithContext(ctx).
		Create(mediaModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, mediaModel.ID)
}

func (r *rMedia) UpdateOne(
	ctx context.Context,
	mediaId uint,
	updateData *entities.MediaUpdate,
) (*entities.Media, error) {
	updates := map[string]interface{}{}

	if updateData.MediaUrl != nil {
		updates["media_url"] = *updateData.MediaUrl
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.UpdatedAt != nil {
		updates["updated_at"] = *updateData.UpdatedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id = ?", mediaId).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, mediaId)
}

func (r *rMedia) DeleteOne(
	ctx context.Context,
	mediaId uint,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.Media{}, mediaId).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *rMedia) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Media, error) {
	media := &models.Media{}

	if res := r.db.WithContext(ctx).
		Model(media).
		Where(query, args...).
		First(media); res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, media.ID)
}

func (r *rMedia) GetMany(
	ctx context.Context,
	query interface{},
	args ...interface{},
) ([]*entities.Media, error) {
	var medias []*models.Media
	if err := r.db.WithContext(ctx).Where(query, args...).Find(&medias).Error; err != nil {
		return nil, err
	}

	var mediaEntities []*entities.Media
	for _, media := range medias {
		mediaEntity := mapper.FromMediaModel(media)
		mediaEntities = append(mediaEntities, mediaEntity)
	}

	return mediaEntities, nil
}