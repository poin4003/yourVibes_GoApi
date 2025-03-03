package implement

import (
	"context"

	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/truncate"

	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentReportRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
)

type sCommentReport struct {
	commentReportRepo commentReportRepo.ICommentReportRepository
	commentRepo       commentReportRepo.ICommentRepository
	notificationRepo  commentReportRepo.INotificationRepository
}

func NewCommentReportImplement(
	commentReportRepo commentReportRepo.ICommentReportRepository,
	commentRepo commentReportRepo.ICommentRepository,
	notificationRepo commentReportRepo.INotificationRepository,
) *sCommentReport {
	return &sCommentReport{
		commentReportRepo: commentReportRepo,
		commentRepo:       commentRepo,
		notificationRepo:  notificationRepo,
	}
}

func (s *sCommentReport) CreateCommentReport(
	ctx context.Context,
	command *commentCommand.CreateReportCommentCommand,
) (result *commentCommand.CreateReportCommentCommandResult, err error) {
	// 1. Check report exist
	commentReportCheck, err := s.commentReportRepo.CheckExist(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		return result, nil
	}

	// 2. Return if report has already exists
	if commentReportCheck {
		return nil, response.NewCustomError(response.ErrCodeCommentReportHasAlreadyExist)
	}

	// 3. Create report
	commentReportEntity, err := commentEntity.NewCommentReport(
		command.UserId,
		command.ReportedCommentId,
		command.Reason,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	commentReport, err := s.commentReportRepo.CreateOne(ctx, commentReportEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 4. Map to result
	return &commentCommand.CreateReportCommentCommandResult{
		CommentReport: mapper.NewCommentReportResult(commentReport),
	}, nil
}

func (s *sCommentReport) HandleCommentReport(
	ctx context.Context,
	command *commentCommand.HandleCommentReportCommand,
) error {
	// 1. Check exist
	commentReportFound, err := s.commentReportRepo.GetById(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if commentReportFound == nil {
		return response.NewDataNotFoundError("comment report not found")
	}

	// 2. Check if report is already handled
	if commentReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &commentEntity.CommentUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	commentUpdated, err := s.commentRepo.UpdateOne(ctx, command.ReportedCommentId, reportedCommentUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Update report status
	commentReportEntity := &commentEntity.CommentReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.commentReportRepo.UpdateMany(ctx, command.ReportedCommentId, commentReportEntity); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Create notification for user
	content := "Your comment has been blocked for violating our policies: " + truncate.TruncateContent(commentUpdated.Content, 20)
	notification, err := notificationEntity.NewNotification(
		commentUpdated.User.FamilyName+" "+commentUpdated.User.Name,
		commentUpdated.User.AvatarUrl,
		commentUpdated.User.ID,
		consts.DEACTIVATE_COMMENT,
		commentUpdated.ID.String(),
		content,
	)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sCommentReport) DeleteCommentReport(
	ctx context.Context,
	command *commentCommand.DeleteCommentReportCommand,
) error {
	// 1. Check exist
	commentReportFound, err := s.commentReportRepo.GetById(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if commentReportFound == nil {
		return response.NewDataNotFoundError("comment report not found")
	}

	// 2. Check if report is already handled
	if commentReportFound.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Delete report
	if err = s.commentReportRepo.DeleteOne(ctx, command.UserId, command.ReportedCommentId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sCommentReport) ActivateComment(
	ctx context.Context,
	command *commentCommand.ActivateCommentCommand,
) error {
	// 1. Check exist
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if commentFound == nil {
		return response.NewDataNotFoundError("comment report not found")
	}

	// 2. Check if comment is already activate
	if commentFound.Status {
		return response.NewCustomError(response.ErrCodeCommentIsAlreadyActivated)
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &commentEntity.CommentUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	commentUpdated, err := s.commentRepo.UpdateOne(ctx, command.CommentId, reportedCommentUpdateEntity)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 4. Delete report
	if err = s.commentReportRepo.DeleteByCommentId(ctx, command.CommentId); err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 5. Create notification for user
	content := "Your comment is activated: " + truncate.TruncateContent(commentUpdated.Content, 20)
	notification, err := notificationEntity.NewNotification(
		commentUpdated.User.FamilyName+" "+commentUpdated.User.Name,
		commentUpdated.User.AvatarUrl,
		commentUpdated.User.ID,
		consts.DEACTIVATE_COMMENT,
		commentUpdated.ID.String(),
		content,
	)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (s *sCommentReport) GetDetailCommentReport(
	ctx context.Context,
	query *commentQuery.GetOneCommentReportQuery,
) (result *commentQuery.CommentReportQueryResult, err error) {
	// 1. Get post report detail
	postReportEntity, err := s.commentReportRepo.GetById(ctx, query.UserId, query.ReportedCommentId)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	if postReportEntity == nil {
		return nil, response.NewDataNotFoundError("post report not found")
	}

	// 2. Map to result
	return &commentQuery.CommentReportQueryResult{
		CommentReport: mapper.NewCommentReportResult(postReportEntity),
	}, nil
}

func (s *sCommentReport) GetManyCommentReport(
	ctx context.Context,
	query *commentQuery.GetManyCommentReportQuery,
) (result *commentQuery.CommentReportQueryListResult, err error) {
	// 1. Get list of comment report
	commentReportEntities, paging, err := s.commentReportRepo.GetMany(ctx, query)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 2. Map to result
	var commentReportResults []*common.CommentReportShortVerResult
	for _, commentReportEntity := range commentReportEntities {
		commentReportResult := mapper.NewCommentReportShortVerResult(commentReportEntity)
		commentReportResults = append(commentReportResults, commentReportResult)
	}

	return &commentQuery.CommentReportQueryListResult{
		CommentReports: commentReportResults,
		PagingResponse: paging,
	}, nil
}
