package repo_impl

import (
	"context"
	"errors"
	"fmt"
	"time"

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
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}
	return mapper.FromMessageModel(&messageModel), nil
}

func (r *rMessage) CreateOne(
	ctx context.Context,
	entity *entities.Message,
) error {
	var messageModel *models.Message
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Conversation{}).
			Where("id = ?", entity.ConversationId).
			Count(&count).Error; err != nil {
		}
		if count <= 0 {
			return response.NewDataNotFoundError("Conversation doesn't exist")
		}
		messageModel = mapper.ToMessageModel(entity)
		if err := tx.Create(messageModel).Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		var userModel *models.User
		if err := tx.Select("id, name").First(&userModel, entity.UserId).Error; err != nil {
			return response.NewDataNotFoundError(err.Error())
		}

		if err := tx.Model(&models.ConversationDetail{}).Where("conversation_id = ?", entity.ConversationId).
			Updates(map[string]interface{}{
				"last_mess":        fmt.Sprintf("%s: %s", userModel.Name, *messageModel.Content),
				"last_mess_status": true,
			}).Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}
		return nil
	})

	return err
}

func (r *rMessage) GetMessagesByConversationId(
	ctx context.Context,
	query *query.GetMessagesByConversationIdQuery,
) ([]*entities.Message, *response.PagingResponse, error) {
	var messageModels []*models.Message
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Message{})

	db = db.Where("conversation_id = ?", query.ConversationId)

	if !query.CreatedAt.IsZero() {
		createAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createAt)
	}

	if query.SortBy != "" {
		shortColumn := ""
		switch query.SortBy {
		case "created_at":
			shortColumn = "created_at"
		}
		if shortColumn != "" {
			if query.IsDescending {
				db = db.Order(shortColumn + " DESC")
			} else {
				db = db.Order(shortColumn)
			}
		}
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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

	if err = db.Offset(offset).Limit(limit).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, family_name, name, avatar_url")
		}).
		Preload("ParentMessage", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, content")
		}).
		Find(&messageModels).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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

func (r *rMessage) DeleteById(
	ctx context.Context,
	id uuid.UUID,
	authenticatedUserId uuid.UUID,
) error {
	message, err := r.GetById(ctx, id)
	if err != nil {
		return response.NewDataNotFoundError(err.Error())
	}

	if authenticatedUserId != message.UserId {
		return response.NewCustomError(response.ErrCantDeleteAnotherMessage)
	}

	res := r.db.WithContext(ctx).Delete(message)
	if res.Error != nil {
		return response.NewServerFailedError(res.Error.Error())
	}
	return nil
}
