package implement

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/pointer"
	"go.uber.org/zap"

	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentValidator "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/validator"
	commentRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sCommentUser struct {
	commentRepo           commentRepo.ICommentRepository
	userRepo              commentRepo.IUserRepository
	postRepo              commentRepo.IPostRepository
	likeUserCommentRepo   commentRepo.ILikeUserCommentRepository
	commentCache          cache.ICommentCache
	postCache             cache.IPostCache
	notificationPublisher *producer.NotificationPublisher
}

func NewCommentUserImplement(
	commentRepo commentRepo.ICommentRepository,
	userRepo commentRepo.IUserRepository,
	postRepo commentRepo.IPostRepository,
	likeUserCommentRepo commentRepo.ILikeUserCommentRepository,
	commentCache cache.ICommentCache,
	postCache cache.IPostCache,
	notificationPublisher *producer.NotificationPublisher,
) *sCommentUser {
	return &sCommentUser{
		commentRepo:           commentRepo,
		userRepo:              userRepo,
		postRepo:              postRepo,
		likeUserCommentRepo:   likeUserCommentRepo,
		commentCache:          commentCache,
		postCache:             postCache,
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

	// 7. Delete cache
	s.commentCache.DeletePostComment(ctx, command.PostId)
	s.postCache.DeletePost(ctx, command.PostId)

	return &commentCommand.CreateCommentResult{
		Comment: mapper.NewCommentResultFromValidateEntity(validateComment),
	}, nil
}

func (s *sCommentUser) UpdateComment(
	ctx context.Context,
	command *commentCommand.UpdateCommentCommand,
) (result *commentCommand.UpdateCommentResult, err error) {
	// 1. Find comment
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return nil, response.NewDataNotFoundError("comment not found")
	}

	// 2. Update comment
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

	// 3. Delete cache
	s.commentCache.DeleteComment(ctx, command.CommentId)

	return &commentCommand.UpdateCommentResult{
		Comment: mapper.NewCommentResultFromEntity(commentUpdate),
	}, nil
}

func (s *sCommentUser) DeleteComment(
	ctx context.Context,
	command *commentCommand.DeleteCommentCommand,
) error {
	// 1. Get post
	comment, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return err
	}

	// 2. Delete comment and child comment in database
	err = s.commentRepo.DeleteCommentAndChildComment(ctx, command.CommentId)
	if err != nil {
		return err
	}

	// 3. Delete cache
	s.commentCache.DeleteComment(ctx, command.CommentId)
	s.commentCache.DeletePostComment(ctx, comment.PostId)
	s.postCache.DeletePost(ctx, comment.PostId)

	return nil
}

func (s *sCommentUser) GetManyComments(
	ctx context.Context,
	query *commentQuery.GetManyCommentQuery,
) (result *commentQuery.GetManyCommentsResult, err error) {
	// 1. Get comment id lít from cache
	commentIDs, paging := s.commentCache.GetPostComment(ctx, query.PostId, query.ParentId, query.Limit, query.Page)

	cacheFailed := false
	if len(commentIDs) == 0 {
		cacheFailed = true
	}

	// 2. Cache hit
	var comments []*commentEntity.Comment
	if !cacheFailed {
		var wg sync.WaitGroup
		var commentMap sync.Map
		cacheErrorOccurred := false

		for _, commentID := range commentIDs {
			wg.Add(1)
			go func(commentID uuid.UUID) {
				defer wg.Done()
				comment := s.commentCache.GetComment(ctx, commentID)
				if comment == nil {
					comment, err = s.commentRepo.GetById(ctx, commentID)
					if err != nil || comment == nil {
						global.Logger.Warn("Failed to get comment", zap.String("commentID", commentID.String()))
						cacheErrorOccurred = true
						s.commentCache.DeleteComment(ctx, commentID)
						s.commentCache.DeletePostComment(ctx, query.PostId)
						return
					}
					s.commentCache.SetComment(ctx, comment)
				}
				commentMap.Store(commentID, comment)
			}(commentID)
		}
		wg.Wait()

		if cacheErrorOccurred {
			cacheFailed = true
		}

		if !cacheFailed {
			for _, commentID := range commentIDs {
				if comment, ok := commentMap.Load(commentID); ok {
					comments = append(comments, comment.(*commentEntity.Comment))
				}
			}
		}
	}

	// 3. Cache miss or handle cache failed
	if cacheFailed {
		global.Logger.Warn("cache failed to get comment, fallback to database")
		var pagingResp *response.PagingResponse
		comments, pagingResp, err = s.commentRepo.GetMany(ctx, query)
		if err != nil {
			return nil, err
		}
		paging = pagingResp

		commentIDs = make([]uuid.UUID, 0, len(comments))
		var wg sync.WaitGroup
		for _, comment := range comments {
			commentIDs = append(commentIDs, comment.ID)
			wg.Add(1)
			go func(c *commentEntity.Comment) {
				defer wg.Done()
				s.commentCache.SetComment(ctx, c)
			}(comment)
		}
		wg.Wait()

		s.commentCache.SetPostComment(ctx, query.PostId, query.ParentId, commentIDs, pagingResp)
	}

	// 4. Get list user like comment
	isLikedListQuery := &commentQuery.CheckUserLikeManyCommentQuery{
		CommentIds:          commentIDs,
		AuthenticatedUserId: query.AuthenticatedUserId,
	}
	isLikedList, err := s.likeUserCommentRepo.CheckUserLikeManyComment(ctx, isLikedListQuery)
	if err != nil {
		return nil, err
	}

	// 5. Map to return
	var commentResults []*common.CommentResultWithLiked
	for _, comment := range comments {
		commentResults = append(commentResults, mapper.NewCommentWithLikedResultFromEntity(comment, isLikedList[comment.ID]))
	}

	return &commentQuery.GetManyCommentsResult{
		Comments:       commentResults,
		PagingResponse: paging,
	}, nil
}
