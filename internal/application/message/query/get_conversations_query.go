package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

type GetConversationsQuery struct {
	UserID uuid.UUID 
}

type GetConversationsQueryResult struct {
	Conversations []models.Conversation 
}

type GetConversationsQueryHandler struct {
	db *gorm.DB 
}

func NewGetConversationsQueryHandler(db *gorm.DB) *GetConversationsQueryHandler {
	return &GetConversationsQueryHandler{db: db}
}

func (h *GetConversationsQueryHandler) Handle(ctx context.Context, query GetConversationsQuery) (*GetConversationsQueryResult, error) {
	// 1. Kiểm tra tính hợp lệ của Query.
	if query.UserID == uuid.Nil {
		return nil, errors.New("user ID is required") 
	}

	// 2. Tìm danh sách các cuộc trò chuyện mà người dùng là thành viên.
	var conversations []models.Conversation
	if err := h.db.WithContext(ctx).
		Joins("JOIN conversation_members ON conversations.id = conversation_members.conversation_id").
		Where("conversation_members.user_id = ?", query.UserID).
		Preload("Members.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Messages.Sender").
		Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err) 
	}

	// 3. Tạo kết quả trả về.
	result := &GetConversationsQueryResult{
		Conversations: conversations,
	}

	return result, nil 
}
