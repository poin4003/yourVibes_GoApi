package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type (
	IConversationRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Conversation, error)
		CreateOne(ctx context.Context, entity *entities.CreateConversation) (*entities.Conversation, error)
		GetManyConversation(ctx context.Context, userId uuid.UUID, query *query.GetManyConversationQuery) ([]*entities.Conversation, *response.PagingResponse, error)
		DeleteById(ctx context.Context, conversationId uuid.UUID, userId uuid.UUID) error
		UpdateOne(ctx context.Context, id uuid.UUID, updateData *entities.ConversationUpdate) (*entities.Conversation, error)
	}
	IMessageRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Message, error)
		CreateOne(ctx context.Context, entity *entities.Message) error
		CreateMany(ctx context.Context, messages []*entities.Message) error
		GetMessagesByConversationId(ctx context.Context, query *query.GetMessagesByConversationIdQuery) ([]*entities.Message, *response.PagingResponse, error)
		DeleteById(ctx context.Context, id uuid.UUID, authenticatedUserId uuid.UUID) error
	}
	IConversationDetailRepository interface {
		GetById(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID) (*entities.ConversationDetail, error)
		CreateOne(ctx context.Context, entity *entities.ConversationDetail) (*entities.ConversationDetail, error)
		GetConversationDetailByConversationId(ctx context.Context, query *query.GetConversationDetailQuery) ([]*entities.ConversationDetail, *response.PagingResponse, error)
		DeleteById(ctx context.Context, userId uuid.UUID, authenticatedUserId uuid.UUID, conversationId uuid.UUID) error
		GetListUserIdByConversationId(ctx context.Context, conversationId uuid.UUID) ([]uuid.UUID, error)
		UpdateOneStatus(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID, updateData *entities.ConversationDetailUpdate) (*entities.ConversationDetail, error)
		CreateMany(ctx context.Context, entities []*entities.ConversationDetail) ([]*entities.ConversationDetail, error)
		TransferOwnerRole(ctx context.Context, userId, authenticatedUserId, conversationId uuid.UUID) error
	}
)
