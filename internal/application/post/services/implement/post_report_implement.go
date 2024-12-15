package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/command"
	post_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sPostReport struct {
	postReportRepo post_report_repo.IPostReportRepository
}

func NewPostReportImplement(
	postReportRepo post_report_repo.IPostReportRepository,
) *sPostReport {
	return &sPostReport{
		postReportRepo: postReportRepo,
	}
}

func (s *sPostReport) CreatePostReport(
	ctx context.Context,
	command *command.CreateReportPostCommand,
) (result *command.CreateReportPostCommandResult, err error) {
	return nil, nil
}
