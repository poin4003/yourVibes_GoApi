package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"go.uber.org/zap"

	"github.com/google/uuid"
	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/validator"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	commentRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sCommentUser struct {
	commentRepo           commentRepo.ICommentRepository
	userRepo              commentRepo.IUserRepository
	postRepo              commentRepo.IPostRepository
	likeUserCommentRepo   commentRepo.ILikeUserCommentRepository
	commentCache          cache.ICommentCache
	notificationPublisher *producer.NotificationPublisher
}

func NewCommentUserImplement(
	commentRepo commentRepo.ICommentRepository,
	userRepo commentRepo.IUserRepository,
	postRepo commentRepo.IPostRepository,
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository,
	commentCache cache.ICommentCache,
	notificationPublisher *producer.NotificationPublisher,
) *sCommentUser {
	return &sCommentUser{
		commentRepo:           commentRepo,
		userRepo:              userRepo,
		postRepo:              postRepo,
		likeUserCommentRepo:   likeUserCommentRepo,
		commentCache:          commentCache,
		notificationPublisher: notificationPublisher,
	}
}

func (s *sCommentUser) CreateComment(
	ctx context.Context,
	command *commentCommand.CreateCommentCommand,
) (result *commentCommand.CreateCommentResult, err error) {
	// 1. Find post
	postFound, err := s.postRepo.GetById(ctx, command.PostId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return nil, response.NewDataNotFoundError("post not found")
	}

	if command.ParentId != nil {
		// 2.1. Get root comment
		var parentComment *commentEntity.Comment
		parentComment, err = s.commentRepo.GetById(ctx, *command.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if parentComment == nil {
			return nil, response.NewDataNotFoundError("parent comment not found")
		}

		// 2. Update rep count +1
		updateComment := &commentEntity.CommentUpdate{
			RepCommentCount: pointer.Ptr(parentComment.RepCommentCount + 1),
		}

		if err = updateComment.ValidateCommentUpdate(); err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		_, err = s.commentRepo.UpdateOne(ctx, parentComment.ID, updateComment)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
	}

	// 4. Create a comment
	newComment, _ := commentEntity.NewComment(
		command.PostId,
		command.UserId,
		command.ParentId,
		command.Content,
	)

	commentCreated, err := s.commentRepo.CreateOne(ctx, newComment)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 5. Update comment count for post
	updatePost := &postEntity.PostUpdate{
		CommentCount: pointer.Ptr(postFound.CommentCount + 1),
	}

	err = updatePost.ValidatePostUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePost)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 4.1.5. Publish to RabbitMQ handle Notification
	notification, _ := notificationEntity.NewNotification(
		commentCreated.User.FamilyName+" "+commentCreated.User.Name,
		commentCreated.User.AvatarUrl,
		postFound.UserId,
		consts.NEW_COMMENT,
		(postFound.ID).String(),
		postFound.Content,
	)

	notifMsg := mapper.NewNotificationResult(notification)
	if err = s.notificationPublisher.PublishNotification(ctx, notifMsg, "notification.single.db_websocket"); err != nil {
		global.Logger.Error("Failed to publish notification result", zap.Error(err))
	}

	// 6. Validate comment after create
	validateComment, err := commentValidator.NewValidatedComment(commentCreated)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &commentCommand.CreateCommentResult{
		Comment: mapper.NewCommentResultFromValidateEntity(validateComment),
	}, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	command *commentCommand.UpdateCommentCommand,
) (result *commentCommand.UpdateCommentResult, err error) {
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return nil, response.NewDataNotFoundError("comment not found")
	}

	updateData := &commentEntity.CommentUpdate{
		Content: command.Content,
	}

	err = updateData.ValidateCommentUpdate()
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	commentUpdate, err := s.commentRepo.UpdateOne(ctx, command.CommentId, updateData)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	return &commentCommand.UpdateCommentResult{
		Comment: mapper.NewCommentResultFromEntity(commentUpdate),
	}, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	command *commentCommand.DeleteCommentCommand,
) error {
	err := s.commentRepo.DeleteCommentAndChildComment(ctx, command.CommentId)
	if err != nil {
		return err
	}

	return nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *commentQuery.GetManyCommentQuery,
) (result *commentQuery.GetManyCommentsResult, err error) {
	result = &commentQuery.GetManyCommentsResult{
		Comments:       nil,
		PagingResponse: nil,
	}

	var queryResult []*commentEntity.Comment
	var paging *response.PagingResponse
	var commentIDs []uuid.UUID
	// Get next layer of comment by root comment
	if query.ParentId != uuid.Nil {
		var parentCommentFound *commentEntity.Comment
		parentCommentFound, err = s.commentRepo.GetById(ctx, query.ParentId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		if parentCommentFound == nil {
			return nil, response.NewDataNotFoundError("parent comment not found")
		}

		queryResult, paging, err = s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return nil, err
		}

		for _, comment := range queryResult {
			commentIDs = append(commentIDs, comment.ID)
		}
	} else {
		// Get first layer if it don't have parent id
		queryResult, paging, err = s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}
		for _, comment := range queryResult {
			commentIDs = append(commentIDs, comment.ID)
		}
	}

	isLikedListQuery := &commentQuery.CheckUserLikeManyCommentQuery{
		CommentIds:          commentIDs,
		AuthenticatedUserId: query.AuthenticatedUserId,
	}

	isLikedList, err := s.likeUserCommentRepo.CheckUserLikeManyComment(ctx, isLikedListQuery)
	if err != nil {
		return nil, err
	}

	var commentResults []*common.CommentResultWithLiked
	for _, comment := range queryResult {
		commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment, isLikedList[comment.ID]))
	}

	result.Comments = commentResults
	result.PagingResponse = paging
	return result, nil
}
