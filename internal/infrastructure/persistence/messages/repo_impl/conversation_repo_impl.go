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

type rConversation struct {
	db *gorm.DB
}

func NewConversationRepositoryImplement(db *gorm.DB) *rConversation {
	return &rConversation{db: db}
}

func (r *rConversation) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Conversation, error) {
	var ConversationModel models.Conversation
	if err := r.db.WithContext(ctx).
		First(&ConversationModel, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}
	return mapper.FromConversationModel(&ConversationModel), nil
}

func (r *rConversation) CreateOne(
	ctx context.Context,
	entity *entities.CreateConversation,
) (*entities.Conversation, error) {
	if len(entity.UserIds) == 2 {
		conversationFound, err := r.findExistingTwoUserConversation(ctx, entity.UserIds[0], entity.UserIds[1])
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
		if conversationFound != nil {
			return nil, response.NewCustomError(response.ErrConversationAlreadyExist, fmt.Sprint(conversationFound.ID))
		}
	}

	var conversationModel *models.Conversation
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		conversationModel = mapper.ToConversationModel(entity)
		if err := tx.Create(conversationModel).Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		for _, userId := range entity.UserIds {
			conversationDetail := &models.ConversationDetail{
				UserId:         userId,
				ConversationId: conversationModel.ID,
				LastMessStatus: true,
				LastMessId:     nil,
			}

			if err := tx.Create(conversationDetail).Error; err != nil {
				return response.NewServerFailedError(err.Error())
			}
		}

		return nil
	})

	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return r.GetById(ctx, conversationModel.ID)
}

func (r *rConversation) GetManyConversation(
	ctx context.Context,
	query *query.GetManyConversationQuery,
) ([]*entities.Conversation, *response.PagingResponse, error) {
	var conversationModels []*models.Conversation
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Conversation{})

	if !query.CreatedAt.IsZero() {
		creatAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", creatAt)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "id":
			sortColumn = "id"
		case "created_at":
			sortColumn = "created_at"

		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
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

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).
		Find(&conversationModels).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var conversationEntities []*entities.Conversation
	for _, conversation := range conversationModels {
		conversationEntities = append(conversationEntities, mapper.FromConversationModel(conversation))
	}

	return conversationEntities, &pagingResponse, nil
}

func (r *rConversation) DeleteById(
	ctx context.Context,
	id uuid.UUID,
) error {
	conversation, err := r.GetById(ctx, id)
	if err != nil {
		return response.NewDataNotFoundError(err.Error())
	}

	res := r.db.WithContext(ctx).Delete(conversation)
	if res.Error != nil {
		return response.NewServerFailedError(res.Error.Error())
	}

	return nil
}

func (r *rConversation) findExistingTwoUserConversation(ctx context.Context, user1, user2 uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	if err := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Joins("JOIN conversation_details cd1 ON conversations.id = cd1.conversation_id").
		Joins("JOIN conversation_details cd2 ON conversations.id = cd2.conversation_id").
		Where("cd1.user_id = ? AND cd2.user_id = ?", user1, user2).
		Where("cd1.user_id != cd2.user_id").
		First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conversation, nil
}
