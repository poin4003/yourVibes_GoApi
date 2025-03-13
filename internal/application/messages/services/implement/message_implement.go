package implement

import (
	"context"

	"github.com/google/uuid"
	messageCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	messageQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	messageEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	messageRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type sMessage struct {
	conversationRepo messageRepo.IConversationRepository
	messageRepo      messageRepo.IMessageRepository
}

func NewMessageImplement(
	conversationRepo messageRepo.IConversationRepository,
	messageRepo messageRepo.IMessageRepository,
) *sMessage {
	return &sMessage{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

func (s *sMessage) GetMessageById(
	ctx context.Context,
	messageId uuid.UUID,
) (result *common.MessageResult, err error) {
	message, err := s.messageRepo.GetById(ctx, messageId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}
	if message == nil {
		return nil, response.NewDataNotFoundError("Message not found")
	}
	return mapper.NewMessageResult(message), nil
}

func (s *sMessage) CreateMessage(
	ctx context.Context,
	command *messageCommand.CreateMessageCommand,
) (result *messageCommand.CreateMessageResult, err error) {
	messageFound, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if messageFound == nil {
		return nil, response.NewDataNotFoundError("Conversation not found")
	}

	if command.ParentId != nil {
		// 2.1. Get root comment
		parentMessage, err := s.messageRepo.GetById(ctx, *command.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if parentMessage == nil {
			return nil, response.NewDataNotFoundError("Parent message not found")
		}
	}

	newMessage, _ := messageEntity.NewMessage(
		command.UserId,
		command.ConversationId,
		command.ParentId,
		&command.Content,
	)

	messageCreated, err := s.messageRepo.CreateOne(ctx, newMessage)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &messageCommand.CreateMessageResult{
		Message: mapper.NewMessageResult(messageCreated),
	}, nil
}

func (s *sMessage) GetMessagesByConversationId(
	ctx context.Context, query *messageQuery.GetMessagesByConversationIdQuery,
) (result *messageQuery.GetMessagesByConversationIdResult, err error) {
	messageEntities, paging, err := s.messageRepo.GetMessagesByConversationId(ctx, query)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	var messageResults []*common.MessageResult
	for _, message := range messageEntities {
		messageResults = append(messageResults, mapper.NewMessageResult(message))
	}

	return &messageQuery.GetMessagesByConversationIdResult{
		Messages:       messageResults,
		PagingResponse: paging,
	}, nil
}
