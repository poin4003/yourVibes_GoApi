package implement

import (
	"context"
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	notificationEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/truncate"
	"net/http"

	commentCommand "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	commentQuery "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	commentEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	commentReportRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
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
	result = &commentCommand.CreateReportCommentCommandResult{
		CommentReport:  nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check report exist
	commentReportCheck, err := s.commentReportRepo.CheckExist(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		return result, nil
	}

	// 2. Return if report has already exists
	if commentReportCheck {
		result.ResultCode = response.ErrCodeCommentReportHasAlreadyExist
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("comment report has already exist")
	}

	// 3. Create report
	commentReportEntity, err := commentEntity.NewCommentReport(
		command.UserId,
		command.ReportedCommentId,
		command.Reason,
	)
	if err != nil {
		return result, err
	}

	commentReport, err := s.commentReportRepo.CreateOne(ctx, commentReportEntity)
	if err != nil {
		return result, err
	}

	// 4. Map to result
	result.CommentReport = mapper.NewCommentReportResult(commentReport)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) HandleCommentReport(
	ctx context.Context,
	command *commentCommand.HandleCommentReportCommand,
) (result *commentCommand.HandleCommentReportCommandResult, err error) {
	result = &commentCommand.HandleCommentReportCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check exist
	commentReportFound, err := s.commentReportRepo.GetById(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if commentReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you dont't need to handle this report again")
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &commentEntity.CommentUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return result, err
	}

	commentUpdated, err := s.commentRepo.UpdateOne(ctx, command.ReportedCommentId, reportedCommentUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("comment report not found")
		}
		return result, err
	}

	// 4. Update report status
	commentReportEntity := &commentEntity.CommentReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.commentReportRepo.UpdateMany(ctx, command.ReportedCommentId, commentReportEntity); err != nil {
		return result, err
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
		return result, err
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) DeleteCommentReport(
	ctx context.Context,
	command *commentCommand.DeleteCommentReportCommand,
) (result *commentCommand.DeleteCommentReportCommandResult, err error) {
	result = &commentCommand.DeleteCommentReportCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check exist
	commentReportFound, err := s.commentReportRepo.GetById(ctx, command.UserId, command.ReportedCommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if report is already handled
	if commentReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you can't delete this report, it's already handled")
	}

	// 3. Delete report
	if err = s.commentReportRepo.DeleteOne(ctx, command.UserId, command.ReportedCommentId); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) ActivateComment(
	ctx context.Context,
	command *commentCommand.ActivateCommentCommand,
) (result *commentCommand.ActivateCommentCommandResult, err error) {
	result = &commentCommand.ActivateCommentCommandResult{
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check exist
	commentFound, err := s.commentRepo.GetById(ctx, command.CommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Check if comment is already activate
	if commentFound.Status {
		result.ResultCode = response.ErrCodeCommentIsAlreadyActivated
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("you dont't need to activate this comment again")
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &commentEntity.CommentUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return result, err
	}

	commentUpdated, err := s.commentRepo.UpdateOne(ctx, command.CommentId, reportedCommentUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("comment report not found")
		}
		return result, err
	}

	// 4. Delete report
	if err = s.commentReportRepo.DeleteByCommentId(ctx, command.CommentId); err != nil {
		return result, err
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
		return result, err
	}

	_, err = s.notificationRepo.CreateOne(ctx, notification)
	if err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) GetDetailCommentReport(
	ctx context.Context,
	query *commentQuery.GetOneCommentReportQuery,
) (result *commentQuery.CommentReportQueryResult, err error) {
	result = &commentQuery.CommentReportQueryResult{
		CommentReport:  nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get post report detail
	postReportEntity, err := s.commentReportRepo.GetById(ctx, query.UserId, query.ReportedCommentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		return result, err
	}

	// 2. Map to result
	result.CommentReport = mapper.NewCommentReportResult(postReportEntity)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) GetManyCommentReport(
	ctx context.Context,
	query *commentQuery.GetManyCommentReportQuery,
) (result *commentQuery.CommentReportQueryListResult, err error) {
	result = &commentQuery.CommentReportQueryListResult{
		CommentReports: nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
		PagingResponse: nil,
	}
	// 1. Get list of comment report
	commentReportEntities, paging, err := s.commentReportRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	// 2. Map to result
	var commentReportResults []*common.CommentReportShortVerResult
	for _, commentReportEntity := range commentReportEntities {
		commentReportResult := mapper.NewCommentReportShortVerResult(commentReportEntity)
		commentReportResults = append(commentReportResults, commentReportResult)
	}

	result.CommentReports = commentReportResults
	result.PagingResponse = paging
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
