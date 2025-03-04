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

// MessageLikeImplementation implements the MessageLikeRepository interface.
type MessageLikeImplementation struct {
	db *gorm.DB
}

// NewMessageLikeImplementation creates a new MessageLikeImplementation.
func NewMessageLikeImplementation(db *gorm.DB) message.MessageLikeRepository {
	return &MessageLikeImplementation{db: db}
}

// Create creates a new message like.
func (r *MessageLikeImplementation) Create(ctx context.Context, messageLike *models.MessageLike) (*models.MessageLike, error) {
	if messageLike.MessageID == uuid.Nil {
		return nil, errors.New("message ID is required")
	}
	if messageLike.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Create(messageLike).Error; err != nil {
		return nil, fmt.Errorf("failed to create message like: %w", err)
	}
	return messageLike, nil
}

// GetByMessageID retrieves message likes by message ID.
func (r *MessageLikeImplementation) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]models.MessageLike, error) {
	if messageID == uuid.Nil {
		return nil, errors.New("message ID is required")
	}
	var messageLikes []models.MessageLike
	if err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Find(&messageLikes).Error; err != nil {
		return nil, fmt.Errorf("failed to get message likes by message ID: %w", err)
	}
	return messageLikes, nil
}

// GetByUserID retrieves message likes by user ID.
func (r *MessageLikeImplementation) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.MessageLike, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	var messageLikes []models.MessageLike
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&messageLikes).Error; err != nil {
		return nil, fmt.Errorf("failed to get message likes by user ID: %w", err)
	}
	return messageLikes, nil
}

// DeleteByMessageID delete all message like of message by message ID.
func (r *MessageLikeImplementation) DeleteByMessageID(ctx context.Context, messageID uuid.UUID) error {
	if messageID == uuid.Nil {
		return errors.New("message ID is required")
	}
	if err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Delete(&models.MessageLike{}).Error; err != nil {
		return fmt.Errorf("failed to delete message likes by message ID: %w", err)
	}
	return nil
}

// DeleteByUserID delete all message like of user by user ID.
func (r *MessageLikeImplementation) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.MessageLike{}).Error; err != nil {
		return fmt.Errorf("failed to delete message likes by user ID: %w", err)
	}
	return nil
}

// Delete delete a message like.
func (r *MessageLikeImplementation) Delete(ctx context.Context, messageLike *models.MessageLike) error {
	if messageLike.MessageID == uuid.Nil {
		return errors.New("message ID is required")
	}
	if messageLike.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if err := r.db.WithContext(ctx).Delete(messageLike).Error; err != nil {
		return fmt.Errorf("failed to delete message like: %w", err)
	}
	return nil
}
