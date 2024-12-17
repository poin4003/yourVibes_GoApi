package implement

import (
	"context"
	"fmt"
	comment_command "github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/mapper"
	comment_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/comment/entities"
	comment_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
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
