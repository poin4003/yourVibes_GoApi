package repo_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	conversationdetail_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/messages/mapper"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/converter"
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
) (*conversationdetail_entity.ConversationDetail, error) {
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

func (r *rConversatioDetail) CreateOne(
	ctx context.Context,
	entity *conversationdetail_entity.ConversationDetail,
) (*conversationdetail_entity.ConversationDetail, error) {
	conversationDetailModel := mapper.ToConversationDetailModel(entity)
	res := r.db.WithContext(ctx).Create(conversationDetailModel)

	if res.Error != nil {
		return nil, response.NewServerFailedError(res.Error.Error())
	}

	return r.GetById(ctx, entity.UserId, entity.ConversationId)
}

func (r *rConversatioDetail) GetConversationDetailByConversationId(
	ctx context.Context,
	query *query.GetConversationDetailQuery,
) ([]*conversationdetail_entity.ConversationDetail, *response.PagingResponse, error) {
	var conversationDetails []*models.ConversationDetail
	var total int64

	db := r.db.WithContext(ctx).Model(&models.ConversationDetail{}).
		Where(" conversation_id = ?", query.ConversationId).
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

	if err := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&conversationDetails).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := response.PagingResponse{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	var conversationDetailEntities []*conversationdetail_entity.ConversationDetail
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
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND conversation_id = ?", userId, conversationId).
		Delete(&conversationdetail_entity.ConversationDetail{})

	if res.Error != nil {
		return response.NewServerFailedError(res.Error.Error())
	}
	return nil
}

func (r *rConversatioDetail) GetListUserIdByConversationId(
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

func (r *rConversatioDetail) UpdateOneStatus(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
	updateData *conversationdetail_entity.ConversationDetailUpdate,
) (*conversationdetail_entity.ConversationDetail, error) {
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

func (r *rConversatioDetail) CreateMany(
	ctx context.Context,
	entities []*conversationdetail_entity.ConversationDetail,
) ([]*conversationdetail_entity.ConversationDetail, error) {
	var conversationDetails []*models.ConversationDetail
	for _, entity := range entities {
		conversationDetails = append(conversationDetails, mapper.ToConversationDetailModel(entity))
	}

	res := r.db.WithContext(ctx).Create(conversationDetails)
	if res.Error != nil {
		return nil, response.NewServerFailedError(res.Error.Error())
	}

	// Lấy thông tin chi tiết cuộc trò chuyện sau khi tạo thành công
	query := &query.GetConversationDetailQuery{
		ConversationId: entities[0].ConversationId, // Sử dụng conversationId của entity đầu tiên
		Page:           1,                          // Trang mặc định, có thể thay đổi tùy vào yêu cầu
		Limit:          10,                         // Giới hạn mặc định, có thể thay đổi
	}
	conversationDetailEntities, _, err := r.GetConversationDetailByConversationId(ctx, query)
	if err != nil {
		return nil, err
	}

	// Trả về các dữ liệu đã tạo và chi tiết cuộc trò chuyện
	return conversationDetailEntities, nil
}
