package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
)

type (
	IConversation interface {
		GetConversationById(ctx context.Context, conversationId uuid.UUID) (result *common.ConversationResult, err error)
		CreateConversation(ctx context.Context, command *command.CreateConversationCommand) (result *command.CreateConversationResult, err error)
		GetManyConversation(ctx context.Context, userId uuid.UUID, query *query.GetManyConversationQuery) (result *query.GetManyConversationQueryResult, err error)
		DeleteConversationById(ctx context.Context, command *command.DeleteConversationCommand) error
		UpdateConversationById(ctx context.Context, command *command.UpdateConversationCommand) (result *command.UpdateConversationCommandResult, err error)
	}
	IMessage interface {
		GetMessageById(ctx context.Context, messageId uuid.UUID) (result *common.MessageResult, err error)
		CreateMessage(ctx context.Context, command *command.CreateMessageCommand) error
		GetMessagesByConversationId(ctx context.Context, query *query.GetMessagesByConversationIdQuery) (result *query.GetMessagesByConversationIdResult, err error)
		DeleteMessageById(ctx context.Context, command *command.DeleteMessageCommand) error
	}
	IConversationDetail interface {
		GetConversationDetailById(ctx context.Context, userId uuid.UUID, conversationId uuid.UUID) (result *common.ConversationDetailResult, err error)
		CreateConversationDetail(ctx context.Context, command *command.CreateConversationDetailCommand) (result *command.CreateConversationDetailResult, err error)
		GetConversationDetailByConversationId(ctx context.Context, query *query.GetConversationDetailQuery) (result *query.GetConversationDetailResult, err error)
		DeleteConversationDetailById(ctx context.Context, command *command.DeleteConversationDetailCommand) error
		UpdateOneStatusConversationDetail(ctx context.Context, command *command.UpdateOneStatusConversationDetailCommand) (err error)
		CreateManyConversationDetail(ctx context.Context, command *command.CreateManyConversationDetailCommand) (result *command.CreateManyConversationDetailResult, err error)
	}
	IMessageMQ interface {
		HandleMessage(ctx context.Context, message *command.CreateMessageCommand) error
	}
)
