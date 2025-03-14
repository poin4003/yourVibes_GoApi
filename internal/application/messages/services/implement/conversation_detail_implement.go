package implement

import (
	"context"

	"github.com/google/uuid"
	conversationDetailCommand "github.com/poin4003/yourVibes_GoApi/internal/application/messages/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/messages/mapper"
	conversationDetailQuery "github.com/poin4003/yourVibes_GoApi/internal/application/messages/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/messages/entities"
	messageRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
)

type sConversationDetail struct {
	conversationRepo       messageRepo.IConversationRepository
	messageRepo            messageRepo.IMessageRepository
	conversationDetailRepo messageRepo.IConversationDetailRepository
}

func NewConversationDetailImplement(
	conversationRepo messageRepo.IConversationRepository,
	messageRepo messageRepo.IMessageRepository,
	conversationDetailRepo messageRepo.IConversationDetailRepository,
) *sConversationDetail {
	return &sConversationDetail{
		conversationRepo:       conversationRepo,
		messageRepo:            messageRepo,
		conversationDetailRepo: conversationDetailRepo,
	}
}

func (s *sConversationDetail) GetConversationDetailById(
	ctx context.Context,
	userId uuid.UUID,
	conversationId uuid.UUID,
) (result *common.ConversationDetailResult, err error) {
	conversationDetail, err := s.conversationDetailRepo.GetById(ctx, userId, conversationId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if conversationDetail == nil {
		return nil, response.NewDataNotFoundError("Conversation detail not found")
	}

	return mapper.NewConversationDetailResult(conversationDetail), nil
}

func (s *sConversationDetail) CreateConversationDetail(
	ctx context.Context,
	command *conversationDetailCommand.CreateConversationDetailCommand,
) (result *conversationDetailCommand.CreateConversationDetailResult, err error) {

	conversation, err := s.conversationRepo.GetById(ctx, command.ConversationId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if conversation == nil {
		return nil, response.NewDataNotFoundError("Conversation not found")
	}

	newconversationDertail, _ := entities.NewConversationDetail(
		command.UserId,
		command.ConversationId,
	)

	conversationCreate, err := s.conversationDetailRepo.CreateOne(ctx, newconversationDertail)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &conversationDetailCommand.CreateConversationDetailResult{
		ConversationDetail: mapper.NewConversationDetailResult(conversationCreate),
	}, nil
}

func (s *sConversationDetail) GetConversationDetailByUsesId(
	ctx context.Context,
	query *conversationDetailQuery.GetConversationDetailQuery,
) (result *conversationDetailQuery.GetConversationDetailResult, err error) {
	conversationDetailEntities, paging, err := s.conversationDetailRepo.GetConversationDetailByUserId(ctx, query)
	if err != nil {
		return result, err
	}

	var conversationDetailResults []*common.ConversationDetailResult
	for _, conversationDetailEntity := range conversationDetailEntities {
		conversationDetailResults = append(conversationDetailResults, mapper.NewConversationDetailResult(conversationDetailEntity))
	}

	return &conversationDetailQuery.GetConversationDetailResult{
		ConversationDetail: conversationDetailResults,
		PagingResponse:     paging,
	}, nil
}
