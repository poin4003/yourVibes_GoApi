package implementations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/message"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"gorm.io/gorm"
)

type MessageUserImplementation struct {
	db *gorm.DB
}

// NewMessageUserImplementation creates a new MessageUserImplementation.
func NewMessageUserImplementation(db *gorm.DB) message.MessageUserRepository {
	return &MessageUserImplementation{db: db}
}

// Create creates a new conversation member.
func (r *MessageUserImplementation) Create(ctx context.Context, conversationMember *models.ConversationMember) (*models.ConversationMember, error) {
	if conversationMember.ConversationID == uuid.Nil {
		return nil, errors.New("conversation ID is required")
	}
	if conversationMember.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Create(conversationMember).Error; err != nil {
		return nil, fmt.Errorf("failed to create conversation member: %w", err)
	}
	return conversationMember, nil
}

// GetByConversationID retrieves conversation member by conversation ID.
func (r *MessageUserImplementation) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]models.ConversationMember, error) {
	if conversationID == uuid.Nil {
		return nil, errors.New("conversation ID is required")
	}
	var conversationMembers []models.ConversationMember
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Find(&conversationMembers).Error; err != nil {
		return nil, fmt.Errorf("failed to get conversation members by conversation ID: %w", err)
	}
	return conversationMembers, nil
}

// GetByUserID retrieves conversation member by User ID.
func (r *MessageUserImplementation) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.ConversationMember, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	var conversationMembers []models.ConversationMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&conversationMembers).Error; err != nil {
		return nil, fmt.Errorf("failed to get conversation members by user ID: %w", err)
	}
	return conversationMembers, nil
}

// DeleteByConversationID delete all conversation member of conversation by conversation ID.
func (r *MessageUserImplementation) DeleteByConversationID(ctx context.Context, conversationID uuid.UUID) error {
	if conversationID == uuid.Nil {
		return errors.New("conversation ID is required")
	}
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Delete(&models.ConversationMember{}).Error; err != nil {
		return fmt.Errorf("failed to delete conversation members by conversation ID: %w", err)
	}
	return nil
}

// DeleteByUserID delete all conversation member of user by user ID.
func (r *MessageUserImplementation) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.ConversationMember{}).Error; err != nil {
		return fmt.Errorf("failed to delete conversation members by user ID: %w", err)
	}
	return nil
}

// Delete delete a conversation member.
func (r *MessageUserImplementation) Delete(ctx context.Context, conversationMember *models.ConversationMember) error {
	if conversationMember.ConversationID == uuid.Nil {
		return errors.New("conversation ID is required")
	}
	if conversationMember.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Delete(conversationMember).Error; err != nil {
		return fmt.Errorf("failed to delete conversation member: %w", err)
	}
	return nil
}
