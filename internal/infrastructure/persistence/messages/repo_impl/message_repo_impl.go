package repo_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/messages/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"gorm.io/gorm"
)

type rMessage struct {
	db *gorm.DB
}

func NewMessageRepositoryImplement(db *gorm.DB) *rMessage {
	return &rMessage{db: db}
}

func (r *rMessage) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Message, error) {
	var messageModel models.Message

	if err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		First(&messageModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.FromMessageModel(&messageModel), nil
}

func (r *rMessage) CreateOne(
	ctx context.Context,
	entity *entities.Message,
) (*entities.Message, error) {
	messageModel := mapper.ToMessageModel(entity)

	res := r.db.WithContext(ctx).Create(messageModel)

	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, messageModel.ID)
}

func (r *rMessage) GetMessagesByConversationId(
	ctx context.Context,
	query *query.GetMessagesByConversationIdQuery,
) ([]*entities.Message, *response.PagingResponse, error) {
	var messageModels []*models.Message
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Message{}).
		Where("conversation_id = ?", query.ConversationId).
		Preload("User").
		Preload("Conversation").
		Order("created_at DESC")

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&messageModels).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var messageEntities []*entities.Message
	for _, messageModel := range messageModels {
		messageEntities = append(messageEntities, mapper.FromMessageModel(messageModel))
	}
	return messageEntities, &pagingResponse, nil
}
