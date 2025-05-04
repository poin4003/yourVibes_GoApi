package repo_impl

import (
	"context"
	"errors"

	"github.com/poin4003/yourVibes_GoApi/internal/consts"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	conversationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/messages/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"
	"gorm.io/gorm"
)

type rConversationDetail struct {
	db *gorm.DB
}

func NewConversationDetailRepositoryImplement(db *gorm.DB) *rConversationDetail {
	return &rConversationDetail{db: db}
}

func (r *rConversationDetail) GetById(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
) (*conversationEntity.ConversationDetail, error) {
	var conversationDetailModel models.ConversationDetail

	if err := r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Preload("User").
		Preload("Conversation").
		First(&conversationDetailModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return mapper.FromConversationDetailModel(&conversationDetailModel), nil
}

func (r *rConversationDetail) CreateOne(
	ctx context.Context,
	entity *conversationEntity.ConversationDetail,
) (*conversationEntity.ConversationDetail, error) {
	conversationDetailModel := mapper.ToConversationDetailModel(entity)
	res := r.db.WithContext(ctx).Create(conversationDetailModel)

	if res.Error != nil {
		return nil, response.NewServerFailedError(res.Error.Error())
	}

	return r.GetById(ctx, entity.UserId, entity.ConversationId)
}

func (r *rConversationDetail) GetConversationDetailByConversationId(
	ctx context.Context,
	query *query.GetConversationDetailQuery,
) ([]*conversationEntity.ConversationDetail, *response.PagingResponse, error) {
	var conversationDetails []*models.ConversationDetail
	var total int64

	db := r.db.WithContext(ctx).Model(&models.ConversationDetail{}).
		Where("conversation_id = ?", query.ConversationId).
		Preload("User").
		Preload("Conversation")

	if err := db.Find(&conversationDetails).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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

	if err = db.WithContext(ctx).Offset(offset).Limit(limit).Find(&conversationDetails).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var conversationDetailEntities []*conversationEntity.ConversationDetail
	for _, conversationDetail := range conversationDetails {
		conversationDetailEntities = append(
			conversationDetailEntities,
			mapper.FromConversationDetailModel(conversationDetail),
		)
	}

	return conversationDetailEntities, &pagingResponse, nil
}

func (r *rConversationDetail) DeleteById(
	ctx context.Context,
	userId uuid.UUID,
	authenticationUserId uuid.UUID,
	conversationId uuid.UUID,
) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check conversation detail exist
		var count int64
		if err := tx.WithContext(ctx).
			Model(&models.ConversationDetail{}).
			Where("conversation_id = ?", conversationId).
			Count(&count).Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}
		if count == 0 {
			return response.NewDataNotFoundError("member or conversation not found")
		}

		if count <= 2 {
			return response.NewCustomError(
				response.ErrCantLeaveConversationIfOnly2Members,
				"Conversation need minimum 2 members",
			)
		}

		// 2. Check role of authenticated userId in conversation
		isOwner, err := r.checkOwnerRole(ctx, conversationId, authenticationUserId)
		if err != nil {
			return err
		}

		// Check Owner of conversation or his own user_id
		if authenticationUserId != userId && !isOwner {
			return response.NewCustomError(
				response.ErrConversationOwnerPermissionRequired,
				"You need to be a owner of this conversation to kick member of group",
			)
		}

		// Check if owner leave conversation
		if isOwner {
			return response.NewCustomError(
				response.ErrCantLeaveConversationIfIsOwners,
				"Owners can't leave conversation",
			)
		}

		// Delete conversation detail
		res := r.db.WithContext(ctx).
			Where("user_id = ? AND conversation_id = ?", userId, conversationId).
			Delete(&conversationEntity.ConversationDetail{})

		if res.Error != nil {
			return response.NewServerFailedError(res.Error.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *rConversationDetail) GetListUserIdByConversationId(
	ctx context.Context,
	conversationId uuid.UUID,
) ([]uuid.UUID, error) {
	var userIds []uuid.UUID

	if err := r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("conversation_id = ?", conversationId).
		Pluck("user_id", &userIds).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	return userIds, nil
}

func (r *rConversationDetail) UpdateOneStatus(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
	updateData *conversationEntity.ConversationDetailUpdate,
) (*conversationEntity.ConversationDetail, error) {
	updates := converter.StructToMap(updateData)
	if len(updates) == 0 {
		return nil, errors.New("no field to update")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Updates(&updates).
		Error; err != nil {

		return nil, response.NewServerFailedError(err.Error())
	}

	return r.GetById(ctx, userId, conversationId)
}

func (r *rConversationDetail) CreateMany(
	ctx context.Context,
	entities []*conversationEntity.ConversationDetail,
) ([]*conversationEntity.ConversationDetail, error) {
	var conversationDetails []*models.ConversationDetail
	var userIds []uuid.UUID
	var conversationIds []uuid.UUID

	for _, entity := range entities {
		conversationDetails = append(conversationDetails, mapper.ToConversationDetailModel(entity))
		userIds = append(userIds, entity.UserId)
		conversationIds = append(conversationIds, entity.ConversationId)
	}

	if err := r.db.WithContext(ctx).Create(&conversationDetails).Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	var createdModels []models.ConversationDetail
	if err := r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("user_id IN ? AND conversation_id IN ?", userIds, conversationIds).
		Preload("User").
		Preload("Conversation").
		Find(&createdModels).Error; err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	var result []*conversationEntity.ConversationDetail
	for _, model := range createdModels {
		result = append(result, mapper.FromConversationDetailModel(&model))
	}

	return result, nil
}

func (r *rConversationDetail) TransferOwnerRole(
	ctx context.Context,
	userId, authenticatedUserId, conversationId uuid.UUID,
) error {
	// 1. Check conversation exists
	isConversationExist, err := r.checkConversationExist(ctx, conversationId)
	if err != nil {
		return err
	}
	if !isConversationExist {
		return response.NewDataNotFoundError("conversation does not exist")
	}
	// 2. Check userId and authenticatedUserId is a member of conversation
	isUserIsMemberOfConversation, err := r.checkMemberExistInConversation(ctx, userId, conversationId)
	if err != nil {
		return err
	}
	if !isUserIsMemberOfConversation {
		return response.NewDataNotFoundError("user is not a member of conversation")
	}

	isAuthenticatedUserIsMemberOfConversation, err := r.checkMemberExistInConversation(ctx, authenticatedUserId, conversationId)
	if err != nil {
		return err
	}
	if !isAuthenticatedUserIsMemberOfConversation {
		return response.NewDataNotFoundError("authenticatedUser is not a member of conversation")
	}

	// 3. Check authenticatedUserId is owner
	isAuthenticatedUserIsOwner, err := r.checkOwnerRole(ctx, conversationId, authenticatedUserId)
	if err != nil {
		return err
	}
	if !isAuthenticatedUserIsOwner {
		return response.NewCustomError(
			response.ErrConversationOwnerPermissionRequired,
			"You need to be a owner of conversation transfer owner role",
		)
	}

	// 4. Transaction to transfer owner role userId <-> authenticatedUserId (make sure conversation only have 1 owner)
	if err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&models.ConversationDetail{}).
			Where("conversation_id = ? AND user_id = ?", conversationId, authenticatedUserId).
			Update("conversation_role", consts.CONVERSATION_MEMBER).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		if err = tx.Model(&models.ConversationDetail{}).
			Where("conversation_id = ? AND user_id = ?", conversationId, userId).
			Update("conversation_role", consts.CONVERSATION_OWNER).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *rConversationDetail) checkOwnerRole(
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

func (r *rConversationDetail) checkConversationExist(
	ctx context.Context,
	conversationId uuid.UUID,
) (isExist bool, err error) {
	var count int64
	if err = r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", conversationId).
		Count(&count).
		Error; err != nil {
	}

	return count > 0, nil
}

func (r *rConversationDetail) checkMemberExistInConversation(
	ctx context.Context,
	userId, conversationId uuid.UUID,
) (isExist bool, err error) {
	var count int64
	if err = r.db.WithContext(ctx).
		Model(&models.ConversationDetail{}).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Count(&count).
		Error; err != nil {
	}

	return count > 0, nil
}
