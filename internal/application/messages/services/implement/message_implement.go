package implement

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/producer"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	messageCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	messageQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	messageEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	messageRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type sMessage struct {
	messageRepo      messageRepo.IMessageRepository
	messagePublisher *producer.MessagePublisher
}

func NewMessageImplement(
	messageRepo messageRepo.IMessageRepository,
	messagePublisher *producer.MessagePublisher,
) *sMessage {
	return &sMessage{
		messageRepo:      messageRepo,
		messagePublisher: messagePublisher,
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
		return nil, err
	}
	return mapper.NewMessageResult(message), nil
}

func (s *sMessage) CreateMessage(
	ctx context.Context,
	command *messageCommand.CreateMessageCommand,
) error {
	// 1. Publish to rabbitmq (push websocket)
	if err := s.messagePublisher.PublishMessage(ctx, command); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 2. Create message in db
	newMessage, _ := messageEntity.NewMessage(
		command.UserId,
		command.ConversationId,
		command.ParentId,
		&command.Content,
	)

	err := s.messageRepo.CreateOne(ctx, newMessage)
	if err != nil {
		return err
	}

	return nil
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

func (s *sMessage) DeleteMessageById(
	ctx context.Context,
	command *command.DeleteMessageCommand,
) error {
	messageFound, err := s.messageRepo.GetById(ctx, *command.MessageId)
	if err != nil {
		return err
	}

	if messageFound == nil {
		return err
	}

	if err := s.messageRepo.DeleteById(ctx, *command.MessageId); err != nil {
		return err
	}

	return nil
}
