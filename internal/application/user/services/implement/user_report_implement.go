package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/command"
	user_report_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sUserReport struct {
	userReportRepo user_report_repo.IUserReportRepository
}

func NewUserReportImplement(
	userReportRepo user_report_repo.IUserReportRepository,
) *sUserReport {
	return &sUserReport{
		userReportRepo: userReportRepo,
	}
}

func (s *sUserReport) CreateUserReport(
	ctx context.Context,
	command *command.CreateReportUserCommand,
) (result *command.CreateReportUserCommandResult, err error) {
	return nil, nil
}
