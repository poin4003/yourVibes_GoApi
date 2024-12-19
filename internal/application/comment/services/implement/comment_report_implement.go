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
	"gorm.io/gorm"
	"net/http"
)

type sCommentReport struct {
	commentReportRepo comment_report_repo.ICommentReportRepository
}

func NewCommentReportImplement(
	commentReportRepo comment_report_repo.ICommentReportRepository,
) *sCommentReport {
	return &sCommentReport{
		commentReportRepo: commentReportRepo,
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
	return nil, nil
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
