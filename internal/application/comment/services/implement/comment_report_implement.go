package implement

import (
	"context"
	"errors"
	"fmt"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	comment_query "github.com/poin4003/yourVibes_GoApi/internal/application/comment/query"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	comment_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sCommentReport struct {
	commentReportRepo comment_report_repo.ICommentReportRepository
	commentRepo       comment_report_repo.ICommentRepository
}

func NewCommentReportImplement(
	commentReportRepo comment_report_repo.ICommentReportRepository,
	commentRepo comment_report_repo.ICommentRepository,
) *sCommentReport {
	return &sCommentReport{
		commentReportRepo: commentReportRepo,
		commentRepo:       commentRepo,
	}
}

func (s *sCommentReport) CreateCommentReport(
	ctx context.Context,
	command *comment_command.CreateReportCommentCommand,
) (result *comment_command.CreateReportCommentCommandResult, err error) {
	result = &comment_command.CreateReportCommentCommandResult{}
	result.CommentReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
	commentReportEntity, err := comment_entity.NewCommentReport(
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
	command *comment_command.HandleCommentReportCommand,
) (result *comment_command.HandleCommentReportCommandResult, err error) {
	result = &comment_command.HandleCommentReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
	if !commentReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("Your dont't need to handle this report again")
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &comment_entity.CommentUpdate{
		Status: pointer.Ptr(false),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return result, err
	}

	_, err = s.commentRepo.UpdateOne(ctx, command.ReportedCommentId, reportedCommentUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("comment report not found")
		}
		return result, err
	}

	// 4. Update report status
	commentReportEntity := &comment_entity.CommentReportUpdate{
		AdminId: pointer.Ptr(command.AdminId),
		Status:  pointer.Ptr(true),
	}

	if err = s.commentReportRepo.UpdateMany(ctx, command.ReportedCommentId, commentReportEntity); err != nil {
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) DeleteCommentReport(
	ctx context.Context,
	command *comment_command.DeleteCommentReportCommand,
) (result *comment_command.DeleteCommentReportCommandResult, err error) {
	result = &comment_command.DeleteCommentReportCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
	if !commentReportFound.Status {
		result.ResultCode = response.ErrCodeReportIsAlreadyHandled
		result.HttpStatusCode = http.StatusBadRequest
		return result, fmt.Errorf("You can't delete this report, it's already handled")
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
	command *comment_command.ActivateCommentCommand,
) (result *comment_command.ActivateCommentCommandResult, err error) {
	result = &comment_command.ActivateCommentCommandResult{}
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
		return result, fmt.Errorf("Your dont't need to activate this comment again")
	}

	// 3. Update reported comment status
	reportedCommentUpdateEntity := &comment_entity.CommentUpdate{
		Status: pointer.Ptr(true),
	}

	if err = reportedCommentUpdateEntity.ValidateCommentUpdate(); err != nil {
		return result, err
	}

	_, err = s.commentRepo.UpdateOne(ctx, command.CommentId, reportedCommentUpdateEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, fmt.Errorf("comment report not found")
		}
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sCommentReport) GetDetailCommentReport(
	ctx context.Context,
	query *comment_query.GetOneCommentReportQuery,
) (result *comment_query.CommentReportQueryResult, err error) {
	result = &comment_query.CommentReportQueryResult{}
	result.CommentReport = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
	query *comment_query.GetManyCommentReportQuery,
) (result *comment_query.CommentReportQueryListResult, err error) {
	result = &comment_query.CommentReportQueryListResult{}
	result.CommentReports = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
	result.PagingResponse = nil
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
