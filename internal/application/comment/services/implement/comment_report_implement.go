package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/comment/command"
	comment_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
	command *command.CreateReportCommentCommand,
) (result *command.CreateReportCommentCommandResult, err error) {
	return nil, nil
}
