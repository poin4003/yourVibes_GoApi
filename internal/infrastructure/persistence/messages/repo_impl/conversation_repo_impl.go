package repo_impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/consts"

	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"

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
	if len(entity.ConversationDetail) == 2 {
		conversationFound, err := r.findExistingTwoUserConversation(
			ctx, entity.ConversationDetail[0].UserId, entity.ConversationDetail[1].UserId,
		)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
		if conversationFound != nil {
			return nil, response.NewCustomError(response.ErrConversationAlreadyExist, fmt.Sprint(conversationFound.ID))
		}
	}

	conversationModel := mapper.ToConversationModel(entity)
	if err := r.db.WithContext(ctx).
		Create(&conversationModel).
		Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return r.GetById(ctx, conversationModel.ID)
}

func (r *rConversation) GetManyConversation(
	ctx context.Context,
	userId uuid.UUID,
	query *query.GetManyConversationQuery,
) ([]*entities.Conversation, *response.PagingResponse, error) {
	var conversationModels []*models.Conversation
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Conversation{}).
		Joins("JOIN conversation_details ON conversation_details.conversation_id = conversations.id").
		Where("conversation_details.user_id = ?", userId).
		Preload("ConversationDetail.User")

	// Filter by created_at
	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createdAt)
	}

	// Filter by group name
	if query.Name != "" {
		db = db.Where("name ILIKE ?", "%"+query.Name+"%")
	}

	// Sort by input column
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

	// Count total record
	err := db.Count(&total).Error
	if err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	// Paging
	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Get list conversation
	if err = db.Offset(offset).Limit(limit).Find(&conversationModels).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	// Paging
	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var conversationResponses []*entities.Conversation
	for _, conversation := range conversationModels {
		// Count members in group
		memberCount := len(conversation.ConversationDetail)

		// Get last messages in conversation
		var lastMess *string
		var lastMessStatus bool

		for _, detail := range conversation.ConversationDetail {
			if detail.UserId == userId {
				lastMess = detail.LastMess
				lastMessStatus = detail.LastMessStatus
				break
			}
		}

		if memberCount >= 3 {
			conversationResponses = append(conversationResponses, &entities.Conversation{
				ID:             conversation.ID,
				Name:           conversation.Name,
				Image:          conversation.Image,
				LastMess:       lastMess,
				LastMessStatus: lastMessStatus,
			})
		} else {
			// If is conversation 1-1, get the order person information instead of group conversation
			var otherUser *entities.User
			for _, detail := range conversation.ConversationDetail {
				if detail.UserId != userId {
					otherUser = &entities.User{
						ID:         detail.User.ID,
						Name:       detail.User.Name,
						FamilyName: detail.User.FamilyName,
						AvatarUrl:  detail.User.AvatarUrl,
					}
					break
				}
			}

			if otherUser != nil {
				conversationResponses = append(conversationResponses, &entities.Conversation{
					ID:             conversation.ID,
					Name:           otherUser.Name,
					FamilyName:     otherUser.FamilyName,
					Avatar:         otherUser.AvatarUrl,
					UserID:         &otherUser.ID,
					LastMess:       lastMess,
					LastMessStatus: lastMessStatus,
				})
			}
		}
	}

	return conversationResponses, &pagingResponse, nil
}

func (r *rConversation) DeleteById(
	ctx context.Context,
	conversationId uuid.UUID,
	userId uuid.UUID,
) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check conversation exist
		isConversationExist, totalMembers, err := r.getConversationAndMembers(ctx, conversationId)
		if err != nil {
			return err
		}

		if !isConversationExist {
			return response.NewDataNotFoundError("conversation not found")
		}

		// 2. Check total of conversation
		if totalMembers >= 3 {
			// Check owner of conversation
			var isOwner bool
			isOwner, err = r.checkOwnerRole(ctx, conversationId, userId)
			if err != nil {
				return err
			}

			// 2. Check Owner of conversation
			if !isOwner {
				return response.NewCustomError(
					response.ErrConversationOwnerPermissionRequired,
					"You need to be a owner of this conversation to delete group",
				)
			}
		}

		// Delete all messages in conversation
		if err = tx.WithContext(ctx).
			Where("conversation_id = ?", conversationId).
			Unscoped().
			Delete(&entities.Message{}).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// Delete all conversation details
		if err = tx.WithContext(ctx).
			Where("conversation_id = ?", conversationId).
			Unscoped().
			Delete(&entities.ConversationDetail{}).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// Delete conversation
		if err = tx.WithContext(ctx).
			Where("id = ?", conversationId).
			Unscoped().
			Delete(&entities.Conversation{}).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *rConversation) findExistingTwoUserConversation(ctx context.Context, user1, user2 uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation

	// Truy vấn để tìm conversation có chính xác 2 người và chứa cả user1, user2
	if err := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Joins("JOIN conversation_details cd1 ON conversations.id = cd1.conversation_id").
		Joins("JOIN conversation_details cd2 ON conversations.id = cd2.conversation_id").
		Where("cd1.user_id = ? AND cd2.user_id = ?", user1, user2).
		// Đếm số lượng người tham gia trong conversation
		Where("conversations.id IN (?)",
			r.db.Model(&models.ConversationDetail{}).
				Select("conversation_id").
				Group("conversation_id").
				Having("COUNT(DISTINCT user_id) = 2"),
		).
		First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

func (r *rConversation) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.ConversationUpdate,
) (*entities.Conversation, error) {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return r.GetById(ctx, id)
}

func (r *rConversation) getConversationAndMembers(
	ctx context.Context,
	conversationId uuid.UUID,
) (isExist bool, totalMembers int64, err error) {
	var totalMember int64
	var count int64
	if err = r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", conversationId).
		Count(&count).
		Error; err != nil {
	}

	if err = r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("conversation_id = ?", conversationId).
		Count(&totalMember).
		Error; err != nil {
	}

	return count > 0, totalMember, nil
}

func (r *rConversation) checkOwnerRole(
	ctx context.Context,
	conversationId uuid.UUID,
	userId uuid.UUID,
) (isOwner bool, err error) {
	var conversationDetail *models.ConversationDetail
	if err = r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("conversation_id = ? AND user_id = ?", conversationId, userId).
		First(&conversationDetail).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, response.NewDataNotFoundError(err.Error())
		}
		return false, response.NewServerFailedError(err.Error())
	}

	if conversationDetail.ConversationRole != consts.CONVERSATION_OWNER {
		return false, nil
	}

	return true, nil
}
