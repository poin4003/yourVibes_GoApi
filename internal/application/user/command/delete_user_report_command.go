package command

import "github.com/google/uuid"

type DeleteUserReportCommand struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
}

type DeleteUserReportCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
