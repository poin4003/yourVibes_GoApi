package implement

import (
	"context"

	"github.com/google/uuid"
	conversationCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	conversationQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	conversationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	conversationRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type sConversation struct {
	conversationRepo conversationRepo.IConversationRepository
}

func NewConversationImplement(
	conversationRepo conversationRepo.IConversationRepository,
) *sConversation {
	return &sConversation{
		conversationRepo: conversationRepo,
	}
}

func (s *sConversation) GetConversationById(
	ctx context.Context,
	conversationId uuid.UUID,
) (result *common.ConversationResult, err error) {
	conversation, err := s.conversationRepo.GetById(ctx, conversationId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if conversation == nil {
		return nil, response.NewDataNotFoundError("Conversation not found")
	}

	return mapper.NewConversationResult(conversation), nil
}

func (s *sConversation) CreateConversation(
	ctx context.Context,
	command *conversationCommand.CreateConversationCommand,
) (result *conversationCommand.CreateConversationResult, err error) {
	conversationEntity, err := conversationEntity.NewConversation(command.Name)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	conversation, err := s.conversationRepo.CreateOne(ctx, conversationEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &conversationCommand.CreateConversationResult{
		Conversation: mapper.NewConversationResult(conversation),
	}, nil
}

func (s *sConversation) GetManyConversation(
	ctx context.Context,
	query *conversationQuery.GetManyConversationQuery,
) (result *conversationQuery.GetManyConversationQueryResult, err error) {
	conversationEntities, paging, err := s.conversationRepo.GetManyConversation(ctx, query)
	if err != nil {
		return result, err
	}

	var conversationResults []*common.ConversationResult
	for _, conversationEntity := range conversationEntities {
		conversationResults = append(conversationResults, mapper.NewConversationResult(conversationEntity))
	}

	return &conversationQuery.GetManyConversationQueryResult{
		Conversation:   conversationResults,
		PagingResponse: paging,
	}, nil
}

func (s *sConversation) DeleteConversationById(
	ctx context.Context,
	command *conversationCommand.DeleteConversationCommand) error {

	//1. Find conversation
	conversationFound, err := s.conversationRepo.GetById(ctx, *command.ConversationId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if conversationFound == nil {
		return response.NewDataNotFoundError("Conversation not found")
	}

	//2. Delete conversation
	if err = s.conversationRepo.DeleteById(ctx, *command.ConversationId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}
