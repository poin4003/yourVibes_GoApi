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
		CreateOne(ctx context.Context, entity *entities.Conversation) (*entities.Conversation, error)
		GetManyConversation(ctx context.Context, query *query.GetManyConversationQuery) ([]*entities.Conversation, *response.PagingResponse, error)
		DeleteById(ctx context.Context, id uuid.UUID) error
	}

	IMessageRepository interface {
		GetById(ctx context.Context, id uuid.UUID) (*entities.Message, error)
		CreateOne(ctx context.Context, entity *entities.Message) error
		GetMessagesByConversationId(ctx context.Context, query *query.GetMessagesByConversationIdQuery) ([]*entities.Message, *response.PagingResponse, error)
		DeleteById(ctx context.Context, id uuid.UUID) error
	}
	IConversationDetailRepository interface {
		GetById(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID) (*entities.ConversationDetail, error)
		CreateOne(ctx context.Context, entity *entities.ConversationDetail) (*entities.ConversationDetail, error)
		GetConversationDetailByIdList(ctx context.Context, query *query.GetConversationDetailQuery) ([]*entities.ConversationDetail, *response.PagingResponse, error)
		DeleteById(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID) error
		GetListUserIdByConversationId(ctx context.Context, conversationId uuid.UUID) ([]uuid.UUID, error)
		UpdateOneStatus(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID, updateData *entities.ConversationDetailUpdate) (*entities.ConversationDetail, error)
	}
)

var (
	localConversation       IConversationRepository
	localMessage            IMessageRepository
	localConversationDetail IConversationDetailRepository
)

func Conversation() IConversationRepository {
	if localConversation == nil {
		panic("repository_implement localConversation not found for interface IConversationRepository")
	}
	return localConversation
}

func InitConversationRepository(i IConversationRepository) {
	localConversation = i
}

func Message() IMessageRepository {
	if localMessage == nil {
		panic("repository_implement localMessage not found for interface IMessageRepository")
	}
	return localMessage
}

func InitMessageRepository(i IMessageRepository) {
	localMessage = i
}

func ConversationDetail() IConversationDetailRepository {
	if localConversationDetail == nil {
		panic("repository_implement localConversationDetail not found for interface IConversationDetailRepository")
	}
	return localConversationDetail
}

func InitConversationDetailRepository(i IConversationDetailRepository) {
	localConversationDetail = i
}
