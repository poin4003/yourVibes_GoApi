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

type rConversatioDetail struct {
	db *gorm.DB
}

func NewConversationDetailRepositoryImplement(db *gorm.DB) *rConversatioDetail {
	return &rConversatioDetail{db: db}
}

func (r *rConversatioDetail) GetById(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
) (*entities.ConversationDetail, error) {
	var conversationDetailModel models.ConversationDetail

	if err := r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Preload("User").
		Preload("Conversation").
		First(&conversationDetailModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.FromConversationDetailModel(&conversationDetailModel), nil

}

func (r *rConversatioDetail) CreateOne(
	ctx context.Context,
	entity *entities.ConversationDetail,
) (*entities.ConversationDetail, error) {
	conversationDetailModel := mapper.ToConversationDetailModel(entity)
	res := r.db.WithContext(ctx).Create(conversationDetailModel)

	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetById(ctx, entity.UserId, entity.ConversationId)
}

func (r *rConversatioDetail) GetConversationDetailByUserId(
	ctx context.Context,
	query *query.GetConversationDetailQuery,
) ([]*entities.ConversationDetail, *response.PagingResponse, error) {
	var conversationDetails []*models.ConversationDetail
	var total int64

	db := r.db.WithContext(ctx).Model(&models.ConversationDetail{}).
		Where("user_id = ? OR conversation_id = ?", query.UserId, query.ConversationId).
		Preload("User").
		Preload("Conversation")
	// if query.ConversationId != uuid.Nil { // Nếu có ConversationId -> Lọc theo User + Conversation
	// 	db = db.Where("user_id = ? AND conversation_id = ?", query.UserId, query.ConversationId)
	// } else { // Nếu chỉ có UserId -> Lấy tất cả theo User
	// 	db = db.Where("user_id = ?", query.UserId)
	// }

	if err := db.Find(&conversationDetails).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

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

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&conversationDetails).Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var conversationDetailEntities []*entities.ConversationDetail
	for _, conversationDetail := range conversationDetails {
		conversationDetailEntities = append(conversationDetailEntities, mapper.FromConversationDetailModel(conversationDetail))
	}

	return conversationDetailEntities, &pagingResponse, nil

}
func (r *rConversatioDetail) DeleteById(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
) error {
	// conversationDetail, err := r.GetById(ctx, userId, conversationId)
	// if err != nil {
	// 	return err
	// }

	res := r.db.WithContext(ctx).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Delete(&entities.ConversationDetail{})

	if res.Error != nil {
		return res.Error
	}
	return nil
}
